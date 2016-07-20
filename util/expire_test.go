// podium
// https://github.com/topfreegames/podium
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2016 Top Free Games <backend@tfgco.com>
// Forked from
// https://github.com/dayvson/go-leaderboard
// Copyright © 2013 Maxwell Dayvson da Silva

package util_test

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/topfreegames/podium/util"
)

var _ = Describe("Expires Helper", func() {
	Describe("No timestamps in leaderboard name", func() {
		It("should get null expiration name without expiration", func() {
			exp, err := util.GetExpireAt("my_league")
			Expect(err).NotTo(HaveOccurred())
			Expect(exp).To(BeEquivalentTo(-1))
		})
	})

	Describe("Yearly expiration", func() {
		It("should get expiration for year 2020", func() {
			exp, err := util.GetExpireAt("league_year2020")
			Expect(err).NotTo(HaveOccurred())

			startTime, err := time.Parse("2006", "2020")
			Expect(err).NotTo(HaveOccurred())

			endTime := startTime.AddDate(2, 0, 0)
			Expect(exp).To(BeEquivalentTo(endTime.Unix()))
		})
	})

	Describe("Custom Day expiration", func() {
		It("should get expiration", func() {
			exp, err := util.GetExpireAt("league_from20201010to20201011")
			Expect(err).NotTo(HaveOccurred())

			startTime, err := time.Parse("20060102", "20201010")
			Expect(err).NotTo(HaveOccurred())
			endTime, err := time.Parse("20060102", "20201011")
			Expect(err).NotTo(HaveOccurred())

			ts := endTime.Add(endTime.Sub(startTime)).Unix()
			Expect(exp).To(BeEquivalentTo(ts))
		})

		It("should return error for invalid timestamp", func() {
			exp, err := util.GetExpireAt("league_from20201039to20201011")
			Expect(exp).To(BeEquivalentTo(-1))
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("parsing time \"20201039\": day out of range"))

			exp, err = util.GetExpireAt("league_from20201010to20201139")
			Expect(exp).To(BeEquivalentTo(-1))
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("parsing time \"20201139\": day out of range"))
		})

		It("should return error for negative duration", func() {
			exp, err := util.GetExpireAt("league_from20201011to20201010")
			Expect(exp).To(BeEquivalentTo(-1))
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("has invalid duration -86400"))
		})
	})

	Describe("Unix Timestamp expiration", func() {
		It("should get expiration", func() {
			start := time.Now()
			end := time.Now().Add(time.Hour * 24)
			exp, err := util.GetExpireAt(fmt.Sprintf("league_from%dto%d", start.Unix(), end.Unix()))
			Expect(err).NotTo(HaveOccurred())

			ts := end.Add(end.Sub(start)).Unix()
			Expect(exp).To(BeEquivalentTo(ts))
		})

		It("should get invalid expiration if timestamps reversed", func() {
			start := time.Now()
			end := time.Now().Add(time.Hour * 24)
			exp, err := util.GetExpireAt(fmt.Sprintf("league_from%dto%d", end.Unix(), start.Unix()))
			Expect(exp).To(BeEquivalentTo(-1))
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("has invalid duration -86400"))
		})
	})

	Describe("Montly expiration", func() {
		It("should get monthly expiration", func() {
			exp, err := util.GetExpireAt("league_year2016month01")
			Expect(err).NotTo(HaveOccurred())

			start, err := time.Parse("20060102", "20160101")
			Expect(err).NotTo(HaveOccurred())

			ts := start.AddDate(0, 2, 0).Unix()
			Expect(exp).To(BeEquivalentTo(ts))
		})

		It("should return error for invalid timestamp", func() {
			exp, err := util.GetExpireAt("league_year2016month99")
			Expect(exp).To(BeEquivalentTo(-1))
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("parsing time \"201699\": month out of range"))
		})
	})

	Describe("Weekly expiration", func() {
		It("should get weekly expiration", func() {
			exp, err := util.GetExpireAt("league_year2016week01")
			Expect(err).NotTo(HaveOccurred())

			end, err := time.Parse("20060102", "20160118")
			Expect(err).NotTo(HaveOccurred())

			ts := end.Unix()
			Expect(exp).To(BeEquivalentTo(ts))
		})
	})

	Describe("Quarter expiration", func() {
		It("should get quarter expiration", func() {
			exp, err := util.GetExpireAt("league_year2016quarter02")
			Expect(err).NotTo(HaveOccurred())

			dummyDate, err := time.Parse("2006", "2016")
			Expect(err).NotTo(HaveOccurred())

			start := dummyDate.AddDate(0, 3, 0)
			end := start.AddDate(0, 6, 0)

			ts := end.Unix()
			Expect(exp).To(BeEquivalentTo(ts))
		})
	})
})
