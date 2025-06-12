package crontask

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"metaLand/app/sync/service/common"
	"metaLand/data/model/crowdfunding"
)

func AddCrowdfundingCornTask(cn *cron.Cron, ctx *common.ServiceContext) {
	fmt.Println("AddCrowdfundingCornTask")
	cn.AddFunc("@every 5s", func() {
		err := LiveCrowdfundingStatusSchedule(ctx.DB)
		if err != nil {
			logx.Errorf("#### update crowdfunding live status: %v\n", err)
		}
	})
	cn.AddFunc("@every 1s", func() {
		err := EndedCrowdfundingStatusSchedule(ctx.DB)
		if err != nil {
			logx.Errorf("#### update crowdfunding ended status: %v\\n", err)
		}
	})
}

func LiveCrowdfundingStatusSchedule(db *gorm.DB) error {
	fmt.Println("LiveCrowdfundingStatusSchedule")
	list, err := crowdfunding.SelectToBeStartedCrowdfundingListWithin1Min(db)
	if err != nil {
		//log.Errorf("#### update crowdfunding live status: %v\n", err)
		return err
	}
	if len(list) > 0 {
		for _, c := range list {
			//log.Infof("#### update crowdfunding live status: %d\n", c.ID)
			logx.Infof("#### update crowdfunding live status: %d\n", c.ID)
			err = crowdfunding.UpdateCrowdfundingStatus(db, c.ID, crowdfunding.Live)
			if err != nil {
				//log.Infof("#### update crowdfunding live status: %v\n", err)
				return err
			}
		}
	}
	return nil
}

func EndedCrowdfundingStatusSchedule(db *gorm.DB) error {
	fmt.Println("EndedCrowdfundingStatusSchedule")
	list, err := crowdfunding.SelectToBeEndedCrowdfundingList(db)
	if err != nil {
		//log.Errorf("#### update crowdfunding ended status: %v\n", err)
		return err
	}
	if len(list) > 0 {
		for _, c := range list {
			err = crowdfunding.UpdateCrowdfundingStatus(db, c.ID, crowdfunding.Ended)
			if err != nil {
				//log.Infof("#### update crowdfunding ended status: %v\n", err)
				return err
			}
		}
	}
	return nil
}
