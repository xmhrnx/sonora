package onnx

import "github.com/mehran9675/sonora/internal/engines/onnx/sherpa"

// piperVoices are the Piper voices sherpa-onnx packages, one download
// each.
var piperVoices = []struct {
	ID, Name, Language, Quality, Size string
}{
	{"en_US-lessac-medium", "Lessac", "English (US)", "Medium", "67 MB"},
	{"en_US-lessac-high", "Lessac", "English (US)", "High", "116 MB"},
	{"en_US-amy-low", "Amy", "English (US)", "Low", "67 MB"},
	{"en_US-ryan-medium", "Ryan", "English (US)", "Medium", "67 MB"},
	{"en_US-ryan-high", "Ryan", "English (US)", "High", "116 MB"},
	{"en_US-ljspeech-high", "LJSpeech", "English (US)", "High", "116 MB"},
	{"en_US-hfc_female-medium", "HFC Female", "English (US)", "Medium", "67 MB"},
	{"en_US-hfc_male-medium", "HFC Male", "English (US)", "Medium", "67 MB"},
	{"en_US-kristin-medium", "Kristin", "English (US)", "Medium", "67 MB"},
	{"en_US-joe-medium", "Joe", "English (US)", "Medium", "67 MB"},
	{"en_US-john-medium", "John", "English (US)", "Medium", "67 MB"},
	{"en_US-libritts_r-medium", "LibriTTS-R", "English (US)", "Medium", "82 MB"},
	{"en_GB-alba-medium", "Alba", "English (UK)", "Medium", "67 MB"},
	{"en_GB-cori-high", "Cori", "English (UK)", "High", "116 MB"},
	{"en_GB-jenny_dioco-medium", "Jenny", "English (UK)", "Medium", "67 MB"},
	{"en_GB-northern_english_male-medium", "Northern English Male", "English (UK)", "Medium", "67 MB"},
	{"en_GB-semaine-medium", "Semaine", "English (UK)", "Medium", "80 MB"},
	{"en_GB-southern_english_female-low", "Southern English Female", "English (UK)", "Low", "67 MB"},
	{"de_DE-thorsten-medium", "Thorsten", "German", "Medium", "67 MB"},
	{"de_DE-thorsten-high", "Thorsten", "German", "High", "116 MB"},
	{"de_DE-kerstin-low", "Kerstin", "German", "Low", "67 MB"},
	{"de_DE-karlsson-low", "Karlsson", "German", "Low", "67 MB"},
	{"de_DE-eva_k-x_low", "Eva", "German", "Low", "27 MB"},
	{"fr_FR-siwis-medium", "Siwis", "French", "Medium", "67 MB"},
	{"fr_FR-upmc-medium", "UPMC", "French", "Medium", "80 MB"},
	{"fr_FR-tom-medium", "Tom", "French", "Medium", "67 MB"},
	{"fr_FR-gilles-low", "Gilles", "French", "Low", "67 MB"},
	{"es_ES-davefx-medium", "Davefx", "Spanish (Spain)", "Medium", "67 MB"},
	{"es_ES-sharvard-medium", "Sharvard", "Spanish (Spain)", "Medium", "80 MB"},
	{"es_MX-claude-high", "Claude", "Spanish (Mexico)", "High", "67 MB"},
	{"es_MX-ald-medium", "Ald", "Spanish (Mexico)", "Medium", "67 MB"},
	{"es_AR-daniela-high", "Daniela", "Spanish (Argentina)", "High", "116 MB"},
	{"it_IT-paola-medium", "Paola", "Italian", "Medium", "67 MB"},
	{"it_IT-riccardo-x_low", "Riccardo", "Italian", "Low", "26 MB"},
	{"pt_BR-faber-medium", "Faber", "Portuguese (Brazil)", "Medium", "67 MB"},
	{"pt_BR-cadu-medium", "Cadu", "Portuguese (Brazil)", "Medium", "67 MB"},
	{"pt_PT-tugao-medium", "Tugão", "Portuguese (Portugal)", "Medium", "67 MB"},
	{"nl_NL-pim-medium", "Pim", "Dutch", "Medium", "67 MB"},
	{"nl_NL-ronnie-medium", "Ronnie", "Dutch", "Medium", "67 MB"},
	{"nl_BE-nathalie-medium", "Nathalie", "Dutch (Belgium)", "Medium", "67 MB"},
	{"sv_SE-nst-medium", "NST", "Swedish", "Medium", "67 MB"},
	{"sv_SE-lisa-medium", "Lisa", "Swedish", "Medium", "67 MB"},
	{"da_DK-talesyntese-medium", "Talesyntese", "Danish", "Medium", "67 MB"},
	{"no_NO-talesyntese-medium", "Talesyntese", "Norwegian", "Medium", "67 MB"},
	{"fi_FI-harri-medium", "Harri", "Finnish", "Medium", "67 MB"},
	{"is_IS-bui-medium", "Búi", "Icelandic", "Medium", "67 MB"},
	{"pl_PL-gosia-medium", "Gosia", "Polish", "Medium", "67 MB"},
	{"pl_PL-darkman-medium", "Darkman", "Polish", "Medium", "67 MB"},
	{"pl_PL-mc_speech-medium", "MC Speech", "Polish", "Medium", "67 MB"},
	{"cs_CZ-jirka-medium", "Jirka", "Czech", "Medium", "67 MB"},
	{"sk_SK-lili-medium", "Lili", "Slovak", "Medium", "67 MB"},
	{"sl_SI-artur-medium", "Artur", "Slovenian", "Medium", "67 MB"},
	{"hu_HU-anna-medium", "Anna", "Hungarian", "Medium", "67 MB"},
	{"hu_HU-imre-medium", "Imre", "Hungarian", "Medium", "67 MB"},
	{"ro_RO-mihai-medium", "Mihai", "Romanian", "Medium", "67 MB"},
	{"ca_ES-upc_ona-medium", "Ona", "Catalan", "Medium", "67 MB"},
	{"el_GR-rapunzelina-low", "Rapunzelina", "Greek", "Low", "67 MB"},
	{"lv_LV-aivars-medium", "Aivars", "Latvian", "Medium", "67 MB"},
	{"sr_RS-serbski_institut-medium", "Serbski Institut", "Serbian", "Medium", "80 MB"},
	{"ru_RU-irina-medium", "Irina", "Russian", "Medium", "67 MB"},
	{"ru_RU-dmitri-medium", "Dmitri", "Russian", "Medium", "67 MB"},
	{"ru_RU-denis-medium", "Denis", "Russian", "Medium", "67 MB"},
	{"uk_UA-ukrainian_tts-medium", "Ukrainian TTS", "Ukrainian", "Medium", "80 MB"},
	{"uk_UA-lada-x_low", "Lada", "Ukrainian", "Low", "27 MB"},
	{"tr_TR-dfki-medium", "DFKI", "Turkish", "Medium", "67 MB"},
	{"tr_TR-fahrettin-medium", "Fahrettin", "Turkish", "Medium", "67 MB"},
	{"cy_GB-gwryw_gogleddol-medium", "Gwryw Gogleddol", "Welsh", "Medium", "67 MB"},
	{"lb_LU-marylux-medium", "Marylux", "Luxembourgish", "Medium", "67 MB"},
	{"ka_GE-natia-medium", "Natia", "Georgian", "Medium", "67 MB"},
	{"ar_JO-kareem-medium", "Kareem", "Arabic", "Medium", "67 MB"},
	{"fa_IR-amir-medium", "Amir", "Persian", "Medium", "67 MB"},
	{"hi_IN-pratham-medium", "Pratham", "Hindi", "Medium", "67 MB"},
	{"vi_VN-vais1000-medium", "VAIS1000", "Vietnamese", "Medium", "67 MB"},
	{"zh_CN-huayan-medium", "Huayan", "Chinese", "Medium", "67 MB"},
}

func piperModel(voice string) model {
	label, size := voice, "67 MB"
	for _, v := range piperVoices {
		if v.ID == voice {
			label, size = v.Name+" · "+v.Language+" · "+v.Quality+" quality", v.Size
		}
	}
	return model{ID: "vits-piper-" + voice, Label: label, Size: size, Family: sherpa.VITS}
}
