package router

import (
	"strconv"

	"github.com/SteinsElite/pickGinS/internal/transaction"
	"github.com/SteinsElite/pickGinS/internal/vault"
	"github.com/gin-gonic/gin"
)

func validPhase(phase string) bool {
	if phase == vault.Week || phase == vault.Month || phase == vault.Year {
		return true
	}
	return false
}

func validTag(tag string) bool {
	if tag == "" ||
		tag == "deposit" ||
		tag == "profit" ||
		tag == "withdraw" {
		return true
	}
	return false
}

// GetTransaction godoc
// @Summary get the transaction info
// @Produce  json
// @Success 200 "the transaction of the page"
// @Param address path string true "user account address"
// @Param tag query string false "tag of the transaction: {deposit,withdraw,claimProfit}, if not specify, get all the category"
// @Param page query int true "inde of page"
// @Param page_size query int true "size of each page"
// @Router /api/v1/transaction/{address} [get]
func GetTransaction(c *gin.Context) {
	// get all the transaction record for the address
	// if param["address"] is "", we should return all kind of transaction
	userAddr := c.Param("address")
	tag := c.Query("tag")
	if !validTag(tag) {
		c.JSON(400, gin.H{
			"err": "Invalid params",
			"msg": "tags should be one of {deposit, withdaw, claimProfit}",
		})
		return
	}
	page, _ := strconv.ParseInt(c.Query("page"), 10, 64)
	pageSize, _ := strconv.ParseInt(c.Query("page_size"), 10, 64)
	if page == 0 || pageSize == 0 {
		c.JSON(400, gin.H{
			"err": "Missing params",
			"msg": "please specify the page &page_size",
		})
		return
	}
	res := transaction.LoadTxFromDb(page, pageSize, tag, userAddr)
	c.JSON(200, gin.H{
		"transaction": res,
	})
}

func GetAnn(c *gin.Context) {

}

func GetSpecificAnn(c *gin.Context) {

}

// GetVolume godoc
// @Summary get the total volume info
// @Produce  json
// @Success 200 {string} json
// @Param range query string true "the duration to query-{7D,1M,1Y}"
// @Router /api/v1/chart/volume [get]
func GetVolume(c *gin.Context) {
	phase := c.Query("range")
	if !validPhase(phase) {
		c.JSON(400, gin.H{
			"error":   "Invalid params",
			"message": "range should be one of {7D,1M,1Y}",
		})
	}
	values := vault.PhasedVolume(phase)
	c.JSON(200, gin.H{
		"points": values,
	})
}

// GetProfit godoc
// @Summary get the phased profit info
// @Produce  json
// @Success 200 "the profit point"
// @Param range query string true "the duration to query-{7D,1M,1Y}"
// @Router /api/v1/chart/profit [get]
func GetProfit(c *gin.Context) {
	phase := c.Query("range")
	if !validPhase(phase) {
		c.JSON(400, gin.H{
			"error":   "Invalid params",
			"message": "range should be one of {7D,1M,1Y}",
		})
	}
	values := vault.PhasedProfit(phase)
	c.JSON(200, gin.H{
		"points": values,
	})
}

// GetRatio godoc
// @Summary get the ratio info
// @Produce  json
// @Success 200 {string} json
// @Router /api/v1/chart/ratio [get]
func GetRatio(c *gin.Context) {
	values := vault.AssetRatio()
	c.JSON(200, gin.H{
		"ratio": values,
	})
}
