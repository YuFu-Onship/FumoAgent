package config

type LanguageStyle struct {
	Chinese  string
	English  string
	Japanese string
}

type LanguageTable struct {
	SideBar_Home     string
	SideBar_Chat     string
	SideBar_Plugin   string
	SideBar_DarkMode string
	SideBar_Setting  string

	SettingPage_Title                       string
	SettingPage_Voice_Title                 string
	SettingPage_Voice_Desc                  string
	SettingPage_Voice_EnableVoice_Title     string
	SettingPage_Voice_EbableVoice_Btn       string
	SettingPage_Voice_CreateVoice_Title     string
	SettingPage_Voice_CreateVoice_Btn       string
	SettingPage_Voice_CreateVoice_Audition  string
	SettingPage_Voice_CreateVoice_Save      string
	SettingPage_Voice_CreateVoice_Restore   string
	SettingPage_Voice_CreateVoice_Preview   string
	SettingPage_Voice_CreateVoice_VoiceName string

	SettingPage_Voice_CoeType       string
	SettingPage_Voice_CoeEngine     string
	SettingPage_Voice_CoeSpeed      string
	SettingPage_Voice_CoeVolume     string
	SettingPage_Voice_CoePitch      string
	SettingPage_Voice_CoeQuality    string
	SettingPage_Voice_CoeIntonation string
	SettingPage_Voice_CoeAccent     string

	SettingPage_Voice_ChoiceVoice    string
	SettingPage_Voice_Default        string
	SettingPage_Voice_Komeiji_Koishi string
	SettingPage_Voice_Koishi_Komeiji string

	SettingPage_Model_Title             string
	SettingPage_Model_Desc              string
	SettingPage_Model_NewAPI_Title      string
	SettingPage_Model_NewAPI_Btn        string
	SettingPage_Model_NewAPI_Save       string
	SettingPage_Model_NewAPI_Restore    string
	SettingPage_Model_NewAPI_Name       string
	SettingPage_Model_NewAPI_Https      string
	SettingPage_Model_NewAPI_Model      string
	SettingPage_Model_NewAPI_Key        string
	SettingPage_Model_ChoiceAPI_Title   string
	SettingPage_Model_ChoiceAPI_Confirm string
	SettingPage_Model_ChoiceAPI_Delete  string

	SettingPage_Character_Title                        string
	SettingPage_Character_Desc                         string
	SettingPage_Character_NewCharacter_Title           string
	SettingPage_Character_NewCharacter_Btn             string
	SettingPage_Character_NewCharacter_Save            string
	SettingPage_Character_NewCharacter_Restore         string
	SettingPage_Character_NewCharacter_Name            string
	SettingPage_Character_NewCharacter_Prompt          string
	SettingPage_Character_NewCharacter_ChoiceCharacter string

	SettingPage_Custom_Title          string
	SettingPage_Custom_Desc           string
	SettingPage_Custom_ChoiceLanguage string
	SettingPage_Custom_ChoiceColor    string

	SettingPage_Info_Title string
	SettingPage_Info_Desc  string
}

func LanguageChinese() LanguageTable {
	return LanguageTable{
		// 侧边栏
		SideBar_Home:     "首页",
		SideBar_Chat:     "聊天",
		SideBar_Plugin:   "插件",
		SideBar_DarkMode: "深色模式",
		SideBar_Setting:  "设置",

		// 设置页通用
		SettingPage_Title: "设置",

		// 语音设置
		SettingPage_Voice_Title:                 "语音设置",
		SettingPage_Voice_Desc:                  "配置语音合成与播放的相关参数",
		SettingPage_Voice_EnableVoice_Title:     "启用语音",
		SettingPage_Voice_EbableVoice_Btn:       "开启语音",
		SettingPage_Voice_CreateVoice_Title:     "创建新语音",
		SettingPage_Voice_CreateVoice_Btn:       "新增语音",
		SettingPage_Voice_CreateVoice_Audition:  "试听",
		SettingPage_Voice_CreateVoice_Save:      "保存",
		SettingPage_Voice_CreateVoice_Restore:   "重置",
		SettingPage_Voice_CreateVoice_Preview:   "预览",
		SettingPage_Voice_CreateVoice_VoiceName: "名称",

		// 语音参数
		SettingPage_Voice_CoeType:       "类型",
		SettingPage_Voice_CoeEngine:     "引擎",
		SettingPage_Voice_CoeSpeed:      "语速",
		SettingPage_Voice_CoeVolume:     "音量",
		SettingPage_Voice_CoePitch:      "音高",
		SettingPage_Voice_CoeIntonation: "顿音",
		SettingPage_Voice_CoeQuality:    "音质",
		SettingPage_Voice_CoeAccent:     "夹音",

		// 语音角色选择
		SettingPage_Voice_ChoiceVoice:    "选择角色语音",
		SettingPage_Voice_Default:        "默认",
		SettingPage_Voice_Komeiji_Koishi: "古明地恋",
		SettingPage_Voice_Koishi_Komeiji: "古明地觉",

		// 模型/接口设置
		SettingPage_Model_Title:             "模型设置",
		SettingPage_Model_Desc:              "配置 LLM 模型接口与 API 密钥",
		SettingPage_Model_NewAPI_Title:      "添加新 API",
		SettingPage_Model_NewAPI_Btn:        "新建接口",
		SettingPage_Model_NewAPI_Save:       "保存接口",
		SettingPage_Model_NewAPI_Restore:    "重置接口",
		SettingPage_Model_NewAPI_Name:       "接口名称",
		SettingPage_Model_NewAPI_Https:      "代理地址",
		SettingPage_Model_NewAPI_Model:      "模型名称",
		SettingPage_Model_NewAPI_Key:        "API 密钥",
		SettingPage_Model_ChoiceAPI_Title:   "选择当前接口",
		SettingPage_Model_ChoiceAPI_Confirm: "保存",
		SettingPage_Model_ChoiceAPI_Delete:  "删除",

		// 角色设置
		SettingPage_Character_Title:                        "角色设置",
		SettingPage_Character_Desc:                         "定义角色的的人格、语气与行为逻辑",
		SettingPage_Character_NewCharacter_Title:           "创建新角色",
		SettingPage_Character_NewCharacter_Btn:             "新增角色",
		SettingPage_Character_NewCharacter_Save:            "保存角色",
		SettingPage_Character_NewCharacter_Restore:         "还原默认",
		SettingPage_Character_NewCharacter_Name:            "角色姓名",
		SettingPage_Character_NewCharacter_Prompt:          "提示词",
		SettingPage_Character_NewCharacter_ChoiceCharacter: "选择当前角色",

		// 个性化/自定义设置
		SettingPage_Custom_Title:          "个性化",
		SettingPage_Custom_Desc:           "语言,配色",
		SettingPage_Custom_ChoiceLanguage: "语言选择",
		SettingPage_Custom_ChoiceColor:    "主题配色",

		// 关于/信息
		SettingPage_Info_Title: "关于",
		SettingPage_Info_Desc:  "软件版本与开发信息",
	}
}
func LanguageEnglish() LanguageTable {
	return LanguageTable{
		// Sidebar
		SideBar_Home:     "Home",
		SideBar_Chat:     "Chat",
		SideBar_Plugin:   "Plugins",
		SideBar_DarkMode: "Dark Mode",
		SideBar_Setting:  "Settings",

		// Setting Page Common
		SettingPage_Title: "Settings",

		// Voice Settings
		SettingPage_Voice_Title:                 "Voice Settings",
		SettingPage_Voice_Desc:                  "Configure text-to-speech and playback parameters.",
		SettingPage_Voice_EnableVoice_Title:     "Enable Voice",
		SettingPage_Voice_EbableVoice_Btn:       "Turn On Voice",
		SettingPage_Voice_CreateVoice_Title:     "Create New Voice",
		SettingPage_Voice_CreateVoice_Btn:       "Add Voice",
		SettingPage_Voice_CreateVoice_Audition:  "Test",
		SettingPage_Voice_CreateVoice_Save:      "Save",
		SettingPage_Voice_CreateVoice_Restore:   "Reset",
		SettingPage_Voice_CreateVoice_Preview:   "Preview",
		SettingPage_Voice_CreateVoice_VoiceName: "Name",

		// Voice Parameters
		SettingPage_Voice_CoeType:       "Type",
		SettingPage_Voice_CoeEngine:     "Engine",
		SettingPage_Voice_CoeSpeed:      "Speed",
		SettingPage_Voice_CoeVolume:     "Volume",
		SettingPage_Voice_CoePitch:      "Pitch",
		SettingPage_Voice_CoeIntonation: "Intonation",
		SettingPage_Voice_CoeQuality:    "Quality",
		SettingPage_Voice_CoeAccent:     "Accent",

		// Voice Choice
		SettingPage_Voice_ChoiceVoice:    "Select Voice Role",
		SettingPage_Voice_Default:        "Default",
		SettingPage_Voice_Komeiji_Koishi: "Komeiji Koishi",
		SettingPage_Voice_Koishi_Komeiji: "Komeiji Satori",

		// Model Settings
		SettingPage_Model_Title:             "Model Settings",
		SettingPage_Model_Desc:              "Configure LLM interface and API keys.",
		SettingPage_Model_NewAPI_Title:      "Add New API",
		SettingPage_Model_NewAPI_Btn:        "New Interface",
		SettingPage_Model_NewAPI_Save:       "Save API",
		SettingPage_Model_NewAPI_Restore:    "Reset API",
		SettingPage_Model_NewAPI_Name:       "API Name",
		SettingPage_Model_NewAPI_Https:      "Proxy Address",
		SettingPage_Model_NewAPI_Model:      "Model Name",
		SettingPage_Model_NewAPI_Key:        "API Key",
		SettingPage_Model_ChoiceAPI_Title:   "Select Current API",
		SettingPage_Model_ChoiceAPI_Confirm: "Save",
		SettingPage_Model_ChoiceAPI_Delete:  "Delete",

		// Character Settings
		SettingPage_Character_Title:                        "Character Settings",
		SettingPage_Character_Desc:                         "Define personality, tone, and behavior logic.",
		SettingPage_Character_NewCharacter_Title:           "Create New Character",
		SettingPage_Character_NewCharacter_Btn:             "Add Character",
		SettingPage_Character_NewCharacter_Save:            "Save Character",
		SettingPage_Character_NewCharacter_Restore:         "Restore Default",
		SettingPage_Character_NewCharacter_Name:            "Name",
		SettingPage_Character_NewCharacter_Prompt:          "Prompt",
		SettingPage_Character_NewCharacter_ChoiceCharacter: "Select Current Character",

		// Customization
		SettingPage_Custom_Title:          "Personalization",
		SettingPage_Custom_Desc:           "UI language and theme color settings.",
		SettingPage_Custom_ChoiceLanguage: "Language",
		SettingPage_Custom_ChoiceColor:    "Theme Color",

		// Info
		SettingPage_Info_Title: "About",
		SettingPage_Info_Desc:  "Software version and development info.",
	}
}

func LanguageJapanese() LanguageTable {
	return LanguageTable{
		// サイドバー
		SideBar_Home:     "ホーム",
		SideBar_Chat:     "チャット",
		SideBar_Plugin:   "プラグイン",
		SideBar_DarkMode: "夜間モード",
		SideBar_Setting:  "設定",

		// 設定ページ共通
		SettingPage_Title: "設定",

		// 音声設定
		SettingPage_Voice_Title:                 "音声設定",
		SettingPage_Voice_Desc:                  "音声合成と再生に関するパラメータを構成します。",
		SettingPage_Voice_EnableVoice_Title:     "音声を有効にする",
		SettingPage_Voice_EbableVoice_Btn:       "音声をオンにする",
		SettingPage_Voice_CreateVoice_Title:     "新しい音声を作成",
		SettingPage_Voice_CreateVoice_Btn:       "音声を追加",
		SettingPage_Voice_CreateVoice_Audition:  "試聴",
		SettingPage_Voice_CreateVoice_Save:      "保存",
		SettingPage_Voice_CreateVoice_Restore:   "リセット",
		SettingPage_Voice_CreateVoice_Preview:   "プレビュー",
		SettingPage_Voice_CreateVoice_VoiceName: "名前",

		// 音声パラメータ
		SettingPage_Voice_CoeType:       "タイプ",
		SettingPage_Voice_CoeEngine:     "エンジン",
		SettingPage_Voice_CoeSpeed:      "速度",
		SettingPage_Voice_CoeVolume:     "音量",
		SettingPage_Voice_CoePitch:      "ピッチ",
		SettingPage_Voice_CoeIntonation: "イントネーション",
		SettingPage_Voice_CoeQuality:    "音質",
		SettingPage_Voice_CoeAccent:     "アクセント",

		// 音声選択
		SettingPage_Voice_ChoiceVoice:    "ボイスを選択",
		SettingPage_Voice_Default:        "デフォルト",
		SettingPage_Voice_Komeiji_Koishi: "古明地 こいし",
		SettingPage_Voice_Koishi_Komeiji: "古明地 さとり",

		// モデル設定
		SettingPage_Model_Title:             "モデル設定",
		SettingPage_Model_Desc:              "LLMモデルのインターフェースとAPIキーを設定します。",
		SettingPage_Model_NewAPI_Title:      "新しいAPIを追加",
		SettingPage_Model_NewAPI_Btn:        "インターフェース作成",
		SettingPage_Model_NewAPI_Save:       "APIを保存",
		SettingPage_Model_NewAPI_Restore:    "APIをリセット",
		SettingPage_Model_NewAPI_Name:       "API名",
		SettingPage_Model_NewAPI_Https:      "プロキシアドレス",
		SettingPage_Model_NewAPI_Model:      "モデル名",
		SettingPage_Model_NewAPI_Key:        "APIキー",
		SettingPage_Model_ChoiceAPI_Title:   "現在のAPIを選択",
		SettingPage_Model_ChoiceAPI_Confirm: "保存",
		SettingPage_Model_ChoiceAPI_Delete:  "削除",

		// キャラクター設定
		SettingPage_Character_Title:                        "キャラクター設定",
		SettingPage_Character_Desc:                         "性格、口調、行動ロジックを定義します。",
		SettingPage_Character_NewCharacter_Title:           "新規キャラクター作成",
		SettingPage_Character_NewCharacter_Btn:             "キャラクター追加",
		SettingPage_Character_NewCharacter_Save:            "キャラクターを保存",
		SettingPage_Character_NewCharacter_Restore:         "デフォルトに戻す",
		SettingPage_Character_NewCharacter_Name:            "名前",
		SettingPage_Character_NewCharacter_Prompt:          "プロンプト",
		SettingPage_Character_NewCharacter_ChoiceCharacter: "現在のキャラクターを選択",

		// カスタマイズ
		SettingPage_Custom_Title:          "カスタマイズ",
		SettingPage_Custom_Desc:           "言語とテーマカラーの設定。",
		SettingPage_Custom_ChoiceLanguage: "言語選択",
		SettingPage_Custom_ChoiceColor:    "テーマカラー",

		// 情報
		SettingPage_Info_Title: "情報",
		SettingPage_Info_Desc:  "ソフトウェアのバージョンと開発情報。",
	}
}
