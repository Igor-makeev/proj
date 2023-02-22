package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"proj/config"
	"proj/internal/entities/models"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

type Client struct {
	Client *http.Client

	address string
}

func NewClient(cfg *config.Config) *Client {
	client := &http.Client{}
	transport := &http.Transport{}
	transport.MaxIdleConns = 20
	client.Transport = transport
	client.Timeout = time.Second * 2
	return &Client{
		Client:  client,
		address: cfg.AccrualSystemAddress,
	}
}

func (c *Client) DoRequest(orderNumber string, UserID int, out chan models.OrderDTO) {
	url := fmt.Sprint(c.address, "/api/orders/", orderNumber)

	var respData models.OrderDTO
Loop:
	for {
		resp, err := c.Client.Get(url)
		if err != nil {
			logrus.Println(err)
		}
		switch resp.StatusCode {
		case http.StatusNoContent:
			logrus.Print("no")
			respData.UserID = UserID
			respData.Number = orderNumber
			respData.Status = models.StatusInvalid
			out <- respData
			break Loop

		case http.StatusTooManyRequests:
			logrus.Print("time to sleep")
			retry, err := strconv.Atoi(resp.Header.Get("Retry-After"))
			if err != nil {
				logrus.Print(err)
			}
			time.Sleep(time.Second * time.Duration(retry))
			continue

		case http.StatusOK:
			logrus.Print("success")
			defer resp.Body.Close()
			err := json.NewDecoder(resp.Body).Decode(&respData)
			if err != nil {
				logrus.Println(err)
			}
			logrus.Print(respData)
			logrus.Print(respData.Number)
			respData.Number = orderNumber
			respData.UserID = UserID
			out <- respData
			break Loop

		case http.StatusInternalServerError:
			logrus.Print("mistakes")
			continue
		}

	}

}
