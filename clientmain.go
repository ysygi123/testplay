package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"test/service"
)

func main() {
	conn, err := grpc.Dial(":8082", grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	p := service.NewGetggClient(conn)

	resp, err := p.GetProdStock(context.Background(), &service.ProductRequest{ProdId: 1})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v\n", resp)
}
