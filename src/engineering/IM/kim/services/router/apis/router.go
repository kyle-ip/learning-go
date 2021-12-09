package apis

import (
	"fmt"
	"hash/crc32"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/klintcheng/kim"
	"github.com/klintcheng/kim/naming"
	"github.com/klintcheng/kim/services/router/conf"
	"github.com/klintcheng/kim/services/router/ipregion"
	"github.com/klintcheng/kim/wire"
	"github.com/sirupsen/logrus"
)

const DefaultLocation = "中国"

type RouterApi struct {
	Naming   naming.Naming
	IpRegion ipregion.IpRegion
	Config   conf.Router
}

type LookUpResp struct {
	UTC      int64    `json:"utc"`
	Location string   `json:"location"`
	Domains  []string `json:"domains"`
}

func (r *RouterApi) Lookup(c iris.Context) {
	ip := kim.RealIP(c.Request())
	token := c.Params().Get("token")

	// step 1
	var location conf.Country
	ipinfo, err := r.IpRegion.Search(ip)
	if err != nil || ipinfo.Country == "0" {
		location = DefaultLocation
	} else {
		location = conf.Country(ipinfo.Country)
	}

	// step 2
	regionId, ok := r.Config.Mapping[location]
	if !ok {
		c.StopWithError(iris.StatusForbidden, err)
		return
	}

	// step 3
	region, ok := r.Config.Regions[regionId]
	if !ok {
		c.StopWithError(iris.StatusInternalServerError, err)
		return
	}

	// step 4
	idc := selectIdc(token, region)

	// step 5
	gateways, err := r.Naming.Find(wire.SNWGateway, fmt.Sprintf("IDC:%s", idc.ID))
	if err != nil {
		c.StopWithError(iris.StatusInternalServerError, err)
		return
	}

	// step 6
	hits := selectGateways(token, gateways, 3)
	domains := make([]string, len(hits))
	for i, h := range hits {
		domains[i] = h.GetMeta()["domain"]
	}

	logrus.WithFields(logrus.Fields{
		"country":  location,
		"regionId": regionId,
		"idc":      idc.ID,
	}).Infof("lookup domain %v", domains)

	_, _ = c.JSON(LookUpResp{
		UTC:      time.Now().Unix(),
		Location: string(location),
		Domains:  domains,
	})
}

func selectIdc(token string, region *conf.Region) *conf.IDC {
	slot := hashcode(token) % len(region.Slots)
	i := region.Slots[slot]
	return &region.Idcs[i]
}

func selectGateways(token string, gateways []kim.ServiceRegistration, num int) []kim.ServiceRegistration {
	if len(gateways) <= num {
		return gateways
	}
	slots := make([]int, 0, len(gateways)*10)
	for i := range gateways {
		for j := 0; j < 10; j++ {
			slots = append(slots, i)
		}
	}
	slot := hashcode(token) % len(slots)
	i := slots[slot]
	res := make([]kim.ServiceRegistration, 0, num)
	for len(res) < num {
		res = append(res, gateways[i])
		i++
		if i >= len(gateways) {
			i = 0
		}
	}
	return res
}

func hashcode(key string) int {
	hash32 := crc32.NewIEEE()
	hash32.Write([]byte(key))
	return int(hash32.Sum32())
}
