// podium
// https://github.com/topfreegames/podium
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2016 Top Free Games <backend@tfgco.com>
// Forked from
// https://github.com/dayvson/go-leaderboard
// Copyright © 2013 Maxwell Dayvson da Silva

package leaderboard_test

import (
	"strconv"

	. "github.com/topfreegames/podium/leaderboard"
	"github.com/topfreegames/podium/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Leaderboard", func() {

	var redisSettings util.RedisSettings
	var redisClient *util.RedisClient

	BeforeEach(func() {
		redisSettings = util.RedisSettings{
			Host:     "localhost",
			Port:     1234,
			Password: "",
		}

		redisClient = util.GetRedisClient(redisSettings)
		conn := redisClient.GetConnection()
		conn.Do("DEL", "test-leaderboard")
	})

	AfterSuite(func() {
		conn := redisClient.GetConnection()
		conn.Do("DEL", "test-leaderboard")
	})

	It("TestSetUserScore", func() {
		testLeaderboard := NewLeaderboard(redisClient, "test-leaderboard", 10)
		dayvson, err := testLeaderboard.SetUserScore("dayvson", 481516)
		Expect(err).To(BeNil())
		arthur, err := testLeaderboard.SetUserScore("arthur", 1000)
		Expect(err).To(BeNil())
		Expect(dayvson.Rank).To(Equal(1))
		Expect(arthur.Rank).To(Equal(2))
	})

	It("TestTotalMembers", func() {
		testLeaderboard := NewLeaderboard(redisClient, "test-leaderboard", 10)
		for i := 0; i < 10; i++ {
			_, err := testLeaderboard.SetUserScore("member_"+strconv.Itoa(i), 1234*i)
			Expect(err).To(BeNil())
		}
		count, err := testLeaderboard.TotalMembers()
		Expect(err).To(BeNil())
		Expect(count).To(Equal(10))
	})

	It("TestRemoveMember", func() {
		testLeaderboard := NewLeaderboard(redisClient, "test-leaderboard", 10)
		for i := 0; i < 10; i++ {
			_, err := testLeaderboard.SetUserScore("member_"+strconv.Itoa(i), 1234*i)
			Expect(err).To(BeNil())
		}
		Expect(testLeaderboard.TotalMembers()).To(Equal(10))
		testLeaderboard.RemoveMember("member_5")
		Expect(testLeaderboard.TotalMembers()).To(Equal(9))
	})

	It("TestTotalPages", func() {
		testLeaderboard := NewLeaderboard(redisClient, "test-leaderboard", 25)
		for i := 0; i < 101; i++ {
			_, err := testLeaderboard.SetUserScore("member_"+strconv.Itoa(i), 1234*i)
			Expect(err).To(BeNil())
		}
		Expect(testLeaderboard.TotalPages()).To(Equal(5))
	})

	It("TestGetUser", func() {
		friendScore := NewLeaderboard(redisClient, "test-leaderboard", 10)
		dayvson, err := friendScore.SetUserScore("dayvson", 12345)
		Expect(err).To(BeNil())
		felipe, err := friendScore.SetUserScore("felipe", 12344)
		Expect(err).To(BeNil())
		Expect(dayvson.Rank).To(Equal(1))
		Expect(felipe.Rank).To(Equal(2))
		friendScore.SetUserScore("felipe", 12346)
		felipe, err = friendScore.GetMember("felipe")
		Expect(err).To(BeNil())
		dayvson, err = friendScore.GetMember("dayvson")
		Expect(err).To(BeNil())
		Expect(felipe.Rank).To(Equal(1))
		Expect(dayvson.Rank).To(Equal(2))
	})

	It("TestGetAroundMe", func() {
		testLeaderboard := NewLeaderboard(redisClient, "test-leaderboard", 25)
		for i := 0; i < 101; i++ {
			_, err := testLeaderboard.SetUserScore("member_"+strconv.Itoa(i), 1234*i)
			Expect(err).To(BeNil())
		}
		users, err := testLeaderboard.GetAroundMe("member_20")
		Expect(err).To(BeNil())
		firstAroundMe := users[0]
		lastAroundMe := users[testLeaderboard.PageSize-1]
		Expect(len(users)).To(Equal(testLeaderboard.PageSize))
		Expect(firstAroundMe.PublicID).To(Equal("member_31"))
		Expect(lastAroundMe.PublicID).To(Equal("member_7"))
	})

	It("TestGetRank", func() {
		testLeaderboard := NewLeaderboard(redisClient, "test-leaderboard", 25)
		for i := 0; i < 101; i++ {
			_, err := testLeaderboard.SetUserScore("member_"+strconv.Itoa(i), 1234*i)
			Expect(err).To(BeNil())
		}
		testLeaderboard.SetUserScore("member_6", 1000)
		Expect(testLeaderboard.GetRank("member_6")).To(Equal(100))
	})

	It("TestGetLeaders", func() {
		testLeaderboard := NewLeaderboard(redisClient, "test-leaderboard", 25)
		for i := 0; i < 1000; i++ {
			_, err := testLeaderboard.SetUserScore("member_"+strconv.Itoa(i+1), 1234*i)
			Expect(err).To(BeNil())
		}
		users, err := testLeaderboard.GetLeaders(1)
		Expect(err).To(BeNil())

		firstOnPage := users[0]
		lastOnPage := users[len(users)-1]
		Expect(len(users)).To(Equal(testLeaderboard.PageSize))
		Expect(firstOnPage.PublicID).To(Equal("member_1000"))
		Expect(firstOnPage.Rank).To(Equal(1))
		Expect(lastOnPage.PublicID).To(Equal("member_976"))
		Expect(lastOnPage.Rank).To(Equal(25))
	})

	It("should add yearly expiration if league supports it", func() {
		leagueID := "test-leaderboard-year2016"
		friendScore := NewLeaderboard(redisClient, leagueID, 10)
		_, err := friendScore.SetUserScore("dayvson", 12345)
		Expect(err).To(BeNil())

		conn := redisClient.GetConnection()
		result, err := conn.Do("TTL", leagueID)
		Expect(err).NotTo(HaveOccurred())

		exp := result.(int64)
		Expect(err).NotTo(HaveOccurred())
		Expect(exp).To(BeNumerically(">", int64(-1)))
	})
})
