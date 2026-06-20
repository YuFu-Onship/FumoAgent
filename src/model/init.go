package model

import (
	"fmt"
	"math/rand/v2"
	"myapp/src/config"
)

// 随机窗口标题
func RandomWindowTitle() string {
	var desc []string
	lang := config.API_CUSTOM_LANGUAGE_GetCurrentLanguage()
	switch lang {
	case "Japanese":
		desc = []string{
			"最強の妖精！",
			"幻想郷の最強",
			"数学の天才！",
			"FumoFumo",
			"満身創痍",
			"バカ！",
			"⑨",
			"Funky！",
			"少女祈祷中...",
			"少女折寿中...",
		}
	case "English":
		desc = []string{
			"Strongest Fairy!",
			"Strongest in Gensokyo",
			"Math Genius!",
			"FumoFumo",
			"Pichuun~",
			"Baka!",
			"⑨",
			"Funky!",
			"Girl's Praying...",
			"Girl's Fraying...",
		}
	default:
		desc = []string{
			"最强大的妖精！",
			"幻想乡の最强",
			"数学天才！",
			"FumoFumo",
			"满身疮痍",
			"Baka！",
			"⑨",
			"Funky！",
			"少女祈祷中...",
			"少女折寿中...",
		}
	}

	randomIndex := rand.IntN(len(desc))
	return fmt.Sprintf("FumoAgent: %s", desc[randomIndex])
}
