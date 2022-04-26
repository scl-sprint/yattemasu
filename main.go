package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"
	"zemmai-dev/yattemasu/domain/model"
	"zemmai-dev/yattemasu/infra/persistence"
	"zemmai-dev/yattemasu/usecase"

	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"googlemaps.github.io/maps"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ChatStatus string

const (
	Normal      = ChatStatus("Normal")
	AreaSetting = ChatStatus("AreaSetting")
)

type ChatChannel struct {
	UserID      string
	ChatStatus  ChatStatus
	LastUpdated time.Time
}

func DotenvLoad() {
	err := godotenv.Load()
	if err != nil {
		log.Print(err)
	}
}

func MeterFromCoordinates(loc1 model.Location, loc2 model.Location) float64 {

	// reference: https://qiita.com/YasumiYasumi/items/3ed8ea69b85dac055381

	const A = 6371008
	const B = 6371008

	const F = (A - B) / A

	lat1 := loc1.Latitude * math.Pi / 180
	lat2 := loc2.Latitude * math.Pi / 180
	lng1 := loc1.Longitude * math.Pi / 180
	lng2 := loc2.Longitude * math.Pi / 180

	phi1 := math.Atan(B / A * math.Tan(lat1))
	phi2 := math.Atan(B / A * math.Tan(lat2))

	f1 := math.Sin(phi1) * math.Sin(phi2)
	f2 := math.Cos(phi1) * math.Cos(phi2)
	f3 := math.Cos(lng1 - lng2)

	X := math.Acos(f1 + f2*f3)

	f4 := math.Sin(X) - X
	f5 := math.Sin(phi1) + math.Sin(phi2)
	f6 := math.Cos(X/2) * math.Cos(X/2)

	f7 := math.Sin(X) + X
	f8 := math.Sin(phi1) - math.Sin(phi2)
	f9 := math.Sin(X/2) * math.Sin(X/2)

	drho := F / 8 * (f4*f5*f5/f6 - f7*f8*f8/f9)

	meter := A * (X + drho)

	return meter
}

func NewPlaceFlexMessage(place model.Place) *linebot.BubbleContainer {

	var titleRatio int = 1
	var contentRatio int = 4

	container := &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Hero: &linebot.ImageComponent{
			Type:       linebot.FlexComponentTypeImage,
			URL:        place.ImageUrl,
			Size:       linebot.FlexImageSizeTypeFull,
			AspectMode: linebot.FlexImageAspectModeTypeCover,
		},
		Body: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type:        linebot.FlexComponentTypeText,
					Text:        place.Name,
					LineSpacing: "20px",
					Size:        linebot.FlexTextSizeTypeXl,
					Weight:      linebot.FlexTextWeightTypeBold,
				},
				&linebot.BoxComponent{
					Type:          linebot.FlexComponentTypeBox,
					Layout:        linebot.FlexBoxLayoutTypeBaseline,
					Spacing:       linebot.FlexComponentSpacingTypeSm,
					PaddingTop:    linebot.FlexComponentPaddingTypeLg,
					PaddingBottom: linebot.FlexComponentPaddingTypeLg,
					Contents: []linebot.FlexComponent{
						&linebot.TextComponent{
							Type:  linebot.FlexComponentTypeText,
							Text:  "住所",
							Flex:  &titleRatio,
							Color: "#aaaaaa",
						},
						&linebot.TextComponent{
							Type:  linebot.FlexComponentTypeText,
							Text:  place.Address,
							Flex:  &contentRatio,
							Color: "#666666",
							Wrap:  true,
						},
					},
				},
			},
		},
	}
	return container
}

func main() {

	channel := []*ChatChannel{}

	DotenvLoad()

	dsn := "user1:user1-passwd@tcp(127.0.0.1:3306)/test-db?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	db.AutoMigrate(&model.User{}, &model.Group{})
	up := persistence.NewUserPersistence(db)
	uu := usecase.NewUserUsecase(up)

	bot, err := linebot.New(os.Getenv("CHANNEL_SECRET"), os.Getenv("CHANNEL_ACCESS_TOKEN"))

	if err != nil {
		log.Fatal(err)
	}

	c, err := maps.NewClient(maps.WithAPIKey(os.Getenv("GMAP_API_KEY")))

	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		events, err := bot.ParseRequest(r)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
		}

		for _, event := range events {

			user, _ := uu.Create(event.Source.UserID)

			switch event.Type {
			case linebot.EventTypeMessage:
				var chatChannel *ChatChannel = &ChatChannel{UserID: event.Source.UserID, ChatStatus: Normal, LastUpdated: event.Timestamp}

				if length := len(channel); length != 0 {
					for index, ch := range channel {
						if ch.UserID == event.Source.UserID {
							chatChannel = channel[index]
							break
						} else {
							if length == index {
								channel = append(channel, chatChannel)
							}
							break
						}
					}
				} else {
					channel = append(channel, chatChannel)
				}

				var replayMessage linebot.SendingMessage

				switch chatChannel.ChatStatus {
				case Normal:
					switch message := event.Message.(type) {
					case *linebot.TextMessage:

						switch message.Text {
						case "エリアを設定":
							locationAction := linebot.NewLocationAction("位置情報を送信！")
							quickReplyItems := linebot.NewQuickReplyItems(linebot.NewQuickReplyButton("", locationAction))
							replayMessage = linebot.NewTextMessage("位置情報を教えてください！").WithQuickReplies(quickReplyItems)
							chatChannel.ChatStatus = AreaSetting
							log.Print(chatChannel)

							if _, err = bot.ReplyMessage(event.ReplyToken, replayMessage).Do(); err != nil {
								log.Print(err)
							}
						default:

							defaultLocation := model.Location{Longitude: 0, Latitude: 0}

							if user.Location == defaultLocation {
								locationAction := linebot.NewLocationAction("まずは位置情報を登録してください！")
								quickReplyItems := linebot.NewQuickReplyItems(linebot.NewQuickReplyButton("", locationAction))
								replayMessage = linebot.NewTextMessage("位置情報を教えてください！").WithQuickReplies(quickReplyItems)
								chatChannel.ChatStatus = AreaSetting
								log.Print(chatChannel)

								if _, err = bot.ReplyMessage(event.ReplyToken, replayMessage).Do(); err != nil {
									log.Print(err)
								}

								return
							}

							req := &maps.TextSearchRequest{
								Query:    message.Text,
								Radius:   1200,
								Location: &maps.LatLng{Lat: user.Location.Latitude, Lng: user.Location.Longitude},
								Type:     maps.PlaceTypeRestaurant,
								Language: "ja",
								OpenNow:  true,
							}

							res, err := c.TextSearch(context.Background(), req)

							if err != nil {
								log.Fatal(err)
							}
							if len(res.Results) == 0 {
								replayMessage = linebot.NewTextMessage("お店が見つかりませんでした...")
								if _, err = bot.ReplyMessage(event.ReplyToken, replayMessage).Do(); err != nil {
									log.Print(err)
								}

							}

							flexCards := []*linebot.BubbleContainer{}

							var slicedResults []maps.PlacesSearchResult

							if len(res.Results) > 8 {
								slicedResults = res.Results[0:8]
							} else {
								slicedResults = res.Results
							}

							for _, result := range slicedResults[0:9] {
								place := model.PlaceFromSearchResult(result)
								log.Print(place.Address)
								if length := MeterFromCoordinates(user.Location, place.Location); length < 1200 {
									log.Print(length)
									flexCards = append(flexCards, NewPlaceFlexMessage(*place))
								}
							}

							if len(flexCards) == 0 {
								replayMessage = linebot.NewTextMessage("お店が見つかりませんでした...")
								if _, err = bot.ReplyMessage(event.ReplyToken, replayMessage).Do(); err != nil {
									log.Print(err)
								}
							}

							flexCarousel := linebot.CarouselContainer{
								Type:     linebot.FlexContainerTypeCarousel,
								Contents: flexCards,
							}

							message := fmt.Sprintf("%s 他 %d 件のお店が見つかりました！", slicedResults[0].Name, len(res.Results))
							if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message), linebot.NewFlexMessage(message, &flexCarousel)).Do(); err != nil {
								log.Print(err)
							}
						}
					}
				case AreaSetting:

					switch message := event.Message.(type) {
					case *linebot.LocationMessage:
						loc := model.Location{Longitude: message.Longitude, Latitude: message.Latitude}
						user, _ = uu.SetLocation(user, loc)
						replayMessage = linebot.NewTextMessage(fmt.Sprintf("あなたの位置情報を「%f, %f」に設定しました！調べたいお店の情報を入力してください！", user.Location.Latitude, user.Location.Longitude))
						chatChannel.ChatStatus = Normal

						if _, err = bot.ReplyMessage(event.ReplyToken, replayMessage).Do(); err != nil {
							log.Print(err)
						}
					}
				}

			}
		}

	})

	if err := http.ListenAndServe(":9000", nil); err != nil {
		log.Fatal(err)
	}
}
