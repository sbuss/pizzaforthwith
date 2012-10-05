package dominos

import (
    "fmt"
    "testing"
    "time"
)

// Helper function for asserting equality
func assertEqual(t *testing.T, b1, b2 interface{}, s string) {
    if b1 != b2 {
        t.Errorf(s)
    }
}

func TestUnmarshall(t *testing.T) {
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

    response, err := UnmarshalXML([]byte(data))
    if err != nil {
        t.Errorf("error: %v", err)
    }
    assertEqual(t, len(response.Statuses), 0, "# of statuses incorrect.")

    loc, _ := time.LoadLocation("America/New_York")
    assertEqual(t, response.Timestamp.Location().String(), loc.String(),
        fmt.Sprintf("Timestamp in wrong location: %s != %s",
            response.Timestamp.Location(), loc))
    assertEqual(t, response.Version, "1.5", "Version does not match")
    assertEqual(t, response.Query.StoreID, "6228", "Store ID does not match")
}

/*
func TestGet(t *testing.T) {
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
    buildMockGet([]byte(data))
    status, err := dominos.status.LatestStatus('5551234567')
}*/
