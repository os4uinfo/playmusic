package main

import (
        "encoding/json"
        "fmt"
        "io/ioutil"
        "log"
        "net/http"
        "os/exec"
)
// 豆瓣播放地址
const URL = "https://douban.fm/j/v2/playlist?channel=-10&kbps=128&client=s%3Amainsite%7Cy%3A3.0&app_name=radio_website&version=100&type=p&sid=427339&pt=0&pb=128"

type DoubanList struct {
        UserInfo   string `json:"warning"`
        VersionMax int    `json:"version_max"`
        SongDetail []Song `json:"song"`
}

type Song struct {
        AlbumTitle       string `json:"albumtitle"`
        Url              string `json:"url"`
        Title            string `json:"title"`
        IsDoubanPlayable bool   `json:"is_douban_playable"`
}

// 获取歌曲信息
func getSong() (string, error) {
        for {
                client := &http.Client{}
                req, err := http.NewRequest("GET", URL, nil)
                if err != nil {
                        log.Fatal(err)
                }
                req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.89 Safari/537.36")

                resp, err := client.Do(req)
                if err != nil {
                        log.Fatalln(err)
                }
                defer resp.Body.Close()

                robots, err := ioutil.ReadAll(resp.Body)
                var m DoubanList
                if err = json.Unmarshal(robots, &m); err != nil {
                        fmt.Printf("Unmarshal err, %v\n", err)
                        return "", err
                }
                if err != nil {
                        log.Fatal(err)
                }
                if len(m.SongDetail) <= 0 {
                        // return "", nil
                        continue
                }
                Notification(m.SongDetail[0].AlbumTitle, m.SongDetail[0].Title)
                return m.SongDetail[0].Url, nil
        }
}

// 调用系统通知
func Notification(ablumtitle, title string) {
        cmd := exec.Command("/usr/bin/notify-send", ablumtitle, title)
        if err := cmd.Start(); err != nil {
                log.Fatal(err)
        }
}

// 调用系统播放器播放
func playMusic(url string) {
        cmd := exec.Command("/usr/bin/mplayer", url)
        // 运行命令
        if err := cmd.Start(); err != nil {
                log.Fatal(err)
        }
        log.Printf("Waiting for command to finish...")
        err := cmd.Wait()
        log.Printf("Command finished with error: %v", err)
}

func main() {
        // 播放
        for {
                songUrl, err := getSong()
                if err != nil {
                        fmt.Printf("err, %v\n", err)
                }
                fmt.Println(songUrl)
                playMusic(songUrl)
        }
}
