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
			respData.UserID = UserID
			respData.Number = orderNumber
			respData.Status = models.StatusInvalid
			break Loop

		case http.StatusTooManyRequests:
			retryTime, err := strconv.Atoi(resp.Header.Get("Retry-After"))
			if err != nil {
				logrus.Println(err)
			}
			time.Sleep(time.Second * time.Duration(retryTime))
			continue

		case http.StatusOK:
			defer resp.Body.Close()
			err := json.NewDecoder(resp.Body).Decode(&respData)
			if err != nil {
				logrus.Println(err)
			}
			respData.UserID = UserID
			break Loop

		case http.StatusInternalServerError:
			continue
		}
	}
	out <- respData
}
