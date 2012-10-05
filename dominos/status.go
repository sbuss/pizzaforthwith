package dominos

import (
    "encoding/xml"
    "errors"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "net/url"
    "time"
)

type Response struct {
    Version         string
    Query           Query
    Timestamp       time.Time
    Statuses        []OrderStatus
}

type Query struct {
    StoreID     string
    OrderKey    string
}

type OrderStatus struct {
    Description     string
    StartTime       time.Time
    OvenTime        time.Time
    RackTime        time.Time
    RouteTime       time.Time
    DeliveryTime    time.Time
}

/*Unmarshal an XML message from Dominos.

Because Go's XML unmarshaler is too young, it doesn't have a way to define
a time format string. Because of this glaring omission we have to use a
custom XML unmarshaler. Of course, the encoding/xml package doesn't export
the Unmarshaler interface so we're faking it till they make it.
*/
func UnmarshalXML(data []byte) (r Response, err error) {
    envelope := struct {
        XMLName     xml.Name    `xml:"Envelope"`
        Body        struct {
            XMLName     xml.Name    `xml:"Body"`
            Response    struct {
                XMLName         xml.Name        `xml:"GetTrackerDataResponse"`
                Version         string
                Query           struct {
                        StoreID     string
                        OrderKey    string
                    }
                AsOf            string
                Statuses        []struct {
                        XMLName         xml.Name    `xml:"OrderStatus"`
                        Description     string      `xml:"OrderDescription"`
                        StartTime       string
                        OvenTime        string
                        RackTime        string
                        RouteTime       string
                        DeliveryTime    string
                    }   `xml:"OrderStatuses>OrderStatus"`
            }
        }
    }{}
    err = xml.Unmarshal(data, &envelope)
    if err != nil {
        msg := fmt.Sprintf("Could not unmarshall response: %s", err)
        log.Fatal(msg)
        return Response{}, errors.New(msg)
    }
    xmlResponse := envelope.Body.Response
    r.Version = xmlResponse.Version
    r.Timestamp = parseAndFixTime(xmlResponse.AsOf)
    r.Query.StoreID = xmlResponse.Query.StoreID
    r.Query.OrderKey = xmlResponse.Query.OrderKey
    for i, status := range xmlResponse.Statuses {
        r.Statuses[i].StartTime = parseAndFixTime(status.StartTime)
        r.Statuses[i].OvenTime = parseAndFixTime(status.OvenTime)
        r.Statuses[i].RackTime = parseAndFixTime(status.RackTime)
        r.Statuses[i].RouteTime = parseAndFixTime(status.RouteTime)
        r.Statuses[i].DeliveryTime = parseAndFixTime(status.DeliveryTime)
    }
    return r, nil
}

/* Parse a Dominos timestamp, which is implictly in eastern time.

Go will put everything into UTC if there isn't a timestamp in the string, so
after parsing the timestamp this function builds a new time.Time in the
America/New_York location.
*/
func parseAndFixTime(s string) time.Time {
    tNull := time.Time{}
    if s == "" {
        return tNull
    }

    // Dominos always reports in eastern time
    location, err := time.LoadLocation("America/New_York")
    if err != nil {
        msg := "Cannot find location"
        log.Fatal(msg)
    }

    format := "2006-01-02T15:04:05"
    t, err := time.Parse(format, s)
    if err != nil {
        msg := fmt.Sprintf("Could not parse time %s", s)
        log.Fatal(msg)
        return tNull
    }
    y, m, d := t.Date()
    H, M, S := t.Clock()
    t2 := time.Date(y, m, d, H, M, S, t.Nanosecond(), location)
    return t2
}

// Query dominos for information about an Order.
func get(params url.Values) ([]byte, error) {
    baseUrl := "http://trkweb.dominos.com/orderstorage/GetTrackerData"
    encodedValues := params.Encode()
    url := baseUrl + "?" + encodedValues
    response, err := http.Get(url)
    if err != nil {
        msg := fmt.Sprintf("Invalid query: %s", url)
        log.Fatal(msg)
        return nil, errors.New(msg)
    }
    responseTextBytes, err := ioutil.ReadAll(response.Body)
    response.Body.Close()
    if err != nil {
        log.Fatal(err)
        return nil, err
    }
    fmt.Println(string(responseTextBytes))
    return responseTextBytes, nil
}

// Check the status of an order, given the Orderer's phone number.
// If the phone number is invalid an empty response will be returned.
func Status(phoneNumber string) (Response, error) {
    values := url.Values{}
    values.Set("Phone", phoneNumber)
    response, err := get(values)
    if err != nil {
        msg := "Invalid phone number"
        log.Fatal(msg)
        return Response{}, errors.New(msg)
    }

    r, err := UnmarshalXML(response)
    if err != nil {
        msg := fmt.Sprintf("Could not unmarshall response: %s", err)
        log.Fatal(msg)
        return Response{}, errors.New(msg)
    }
    return r, nil
}
