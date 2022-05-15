package os_utils

import (
	"io/ioutil"
	"os"
	"os/user"
	"strconv"
	"syscall"
	"time"
)

type LoginUserStat struct {
	UserStatMap map[string]UserStat
}

type UserStat struct {
	User          string
	Tty           string
	From          string
	LoginDuration int
	Idle          string
	Jcpu          string
	Pcpu          string
	What          string
}

func GetLoginUserStat() (loginUserStat *LoginUserStat, err error) {
	var files []os.FileInfo
	if files, err = ioutil.ReadDir("/dev/pts"); err != nil {
		return
	}
	now := time.Now()

	userStatMap := map[string]UserStat{}
	for _, file := range files {
		stat, ok := file.Sys().(*syscall.Stat_t)
		if !ok {
			continue
		}

		var owner string
		uid := strconv.Itoa(int(stat.Uid))
		u, err := user.LookupId(uid)
		if err != nil {
			owner = uid
		} else {
			owner = u.Username
		}

		loginDuration := now.Unix() - file.ModTime().Unix()

		userStatMap[owner] = UserStat{
			User:          owner,
			Tty:           file.Name(),
			LoginDuration: int(loginDuration),
		}
	}
	loginUserStat = &LoginUserStat{
		UserStatMap: userStatMap,
	}
	return
}
