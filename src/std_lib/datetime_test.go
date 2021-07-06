package std_lib

import (
    "errors"
    "fmt"
    "testing"
    "time"
)

var timeLayout = "2006-01-02 15:04:05"

func TestLocalTimezone(t *testing.T) {
    //tz, ok := syscall.Getenv("TZ")
    //switch {
    //case !ok:
    //   z, err := loadZoneFile("", "/etc/localtime")
    //   if err == nil {
    //       localLoc = *z
    //       localLoc.name = "Local"
    //       return
    //   }
    //case tz != "" && tz != "UTC":
    //   if z, err := loadLocation(tz); err == nil {
    //       localLoc = *z
    //       return
    //   }
    //}
}

func TestTimeParse(t *testing.T) {
    tm, _ := time.Parse(timeLayout, "2016-06-13 09:14:00")
    fmt.Println(time.Now().Sub(tm).Hours())
}

type OftenTime time.Time

func (self OftenTime) MarshalJSON() ([]byte, error) {
    t := time.Time(self)
    if y := t.Year(); y < 0 || y >= 10000 {
        return nil, errors.New("Time.MarshalJSON: year outside of range [0,9999]")
    }
    // 注意 `"2006-01-02 15:04:05"`。因为是 JSON，双引号不能少
    return []byte(t.Format(`"2006-01-02 15:04:05"`)), nil
}

func TestMarshalJSON(t *testing.T) {

}

func TestTimeTruncate(t *testing.T) {
    tm, _ := time.ParseInLocation("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:00:00"), time.Local)
    fmt.Println(tm)

    tm, _ = time.ParseInLocation("2006-01-02 15:04:05", "2016-06-13 15:34:39", time.Local)

    // 整点（向下取整）
    fmt.Println(tm.Truncate(1 * time.Hour))

    // 整点（最接近）
    fmt.Println(tm.Round(1 * time.Hour))

    // 整分（向下取整）
    fmt.Println(tm.Truncate(1 * time.Minute))

    // 整分（最接近）
    fmt.Println(tm.Round(1 * time.Minute))

    t2, _ := time.ParseInLocation("2006-01-02 15:04:05", tm.Format("2006-01-02 15:00:00"), time.Local)
    fmt.Println(t2)
}

func TestTimeout(t *testing.T) {
    c := make(chan int)

    go func() {
        // time.Sleep(1 * time.Second)
        time.Sleep(3 * time.Second)
        <-c
    }()

    select {
    case c <- 1:
        fmt.Println("channel...")
    case <-time.After(2 * time.Second):
        close(c)
        fmt.Println("timeout...")
    }
}

func TestStopWatch(t *testing.T) {
    start := time.Now()
    timer := time.AfterFunc(2*time.Second, func() {
        fmt.Println("after func callback, elaspe:", time.Now().Sub(start))
    })

    time.Sleep(1 * time.Second)
    // time.Sleep(3*time.Second)
    // Reset 在 Timer 还未触发时返回 true；触发了或 Stop 了，返回 false
    if timer.Reset(3 * time.Second) {
        fmt.Println("timer has not trigger!")
    } else {
        fmt.Println("timer had expired or stop!")
    }

    time.Sleep(10 * time.Second)

    // output:
    // timer has not trigger!
    // after func callback, elaspe: 4.00026461s
}
