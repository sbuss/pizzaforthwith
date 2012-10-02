package main

import (
    "fmt"
    "encoding/xml"
)

type Envelope struct {
    XMLName     xml.Name    `xml:"Envelope"`
    Body        Body
}

type Body struct {
    XMLName     xml.Name    `xml:"Body"`
    Response    Response
}

type Response struct {
    XMLName         xml.Name        `xml:"GetTrackerDataResponse"`
    Version         string
    Query           Query
    Timestamp       string          `xml:"AsOf"`
    Statuses        []OrderStatus   `xml:"OrderStatuses>OrderStatus"`
}

type Query struct {
    StoreID     string
    OrderKey    string
}

type OrderStatus struct {
    XMLName         xml.Name    `xml:"OrderStatus"`
    Description     string      `xml:"OrderDescription"`
    StartTime       string
    OvenTime        string
    RackTime        string
    RouteTime       string
    DeliveryTime    string
}

func main() {
    v := Envelope{}

    /*data := `
    <OrderStatuses>
        <OrderStatus>
            <OrderDescription>2 Large Pizzas</OrderDescription>
            <StartTime
        </OrderStatus>
    </OrderStatuses>
    `*/
    data := `
    <soap:Envelope>
        <soap:Body>
            <GetTrackerDataResponse>
                <Version>1.5</Version>
                <Query>
                    <StoreID>6228</StoreID>
                    <OrderKey>622834370420</OrderKey>
                </Query>
                <AsOf>2012-10-02T02:57:08</AsOf>
                <OrderStatuses/>
            </GetTrackerDataResponse>
        </soap:Body>
    </soap:Envelope>
    `

    err := xml.Unmarshal([]byte(data), &v)
    if err != nil {
        fmt.Printf("error: %v", err)
        return
    }
    fmt.Println(v)
    fmt.Printf("# Statuses: %d\n", len(v.Body.Response.Statuses))
    //fmt.Printf("XMLName: %#v\n", v.XMLName)
    //fmt.Printf("Description: %q\n", v.Description)
}
