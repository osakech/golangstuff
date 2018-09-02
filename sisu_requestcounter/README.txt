1. cd into this directory
2. run the tool
    go run main.go
3. open the browser and get to adress "http://localhost:8080/counter/"
    If the tool is working you should see the line:
    "Calls from the last 60 sec -> 7"

    Here you can see the timeframe "60 sec"
    and number of calls within this frame "7"

    Now if you hit the refresh button the number of calls should rise
    and time frame should move with it. For more convenience change the 
    variable "timeframeInSec" in main.go to a lower value
4. Run the Tests
    Last but not least you can also run the tests
    go test -v