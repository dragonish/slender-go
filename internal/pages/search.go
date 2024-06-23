package pages

import (
	"html/template"
	"slender/internal/database"
	"slender/internal/global"
	"slender/internal/icons"
	"slender/internal/model"
)

// generateSearchModule returns search module.
func generateSearchModule(clearSearchTip string, inHomeSearch string, useInHomeSearch string) template.HTML {
	searchModuleTpl := ""

	if global.Config.ShowSearchInput {
		searchEngineList := make([]model.HomeSearchEngineListItem, 0)
		rErr := database.GetHomeSearchEngines(&searchEngineList)
		if rErr != nil {
			return template.HTML(searchModuleTpl)
		}

		searchModuleTpl = `<div id="slender-search" class="slender-module-container slender-top-module"><div class="slender-search-input-module"><input id="slender-search-input" type="text" name="search" autocomplete="off" autofocus /><span class="slender-search-input-prefix"></span><div id="slender-search-input-prefix-overlay"></div><img class="slender-search-input-suffix slender-search-built-in-icon" src="` + icons.GetBuiltInIcon("mdi-searchweb") + `" alt="search" /><div id="slender-search-input-suffix-overlay"></div><label for="slender-search-input" id="slender-clear-button"> >> ` + clearSearchTip + `</label><span id="slender-search-total">&nbsp;</span></div><div id="slender-search-engine-module"><div class="slender-search-engine-dialog"><div class="slender-search-engine-config"><input id="slender-search-enable-in-home" type="checkbox" /><label for="slender-search-enable-in-home" class="slender-search-checkbox-label">` + useInHomeSearch + `</label></div><div id="slender-search-engine-list"><ul><li class="slender-search-engine-item"><a class="slender-search-engine-link" href="javascript:void(0)" data-id="0" data-name="` + inHomeSearch + `" data-url="" data-method="" data-body="">` + renderSearchEngineIcon(inHomeSearch, "mdi-homesearch") + `<span class="slender-search-engine-text">` + inHomeSearch + `</span></a></li>` + renderSearchEngineList(searchEngineList) + `</ul></div></div></div></div>`
	}

	return template.HTML(searchModuleTpl)
}

func renderSearchEngineList(searchEngines []model.HomeSearchEngineListItem) string {
	tpl := ""

	for _, item := range searchEngines {
		tpl += `<li class="slender-search-engine-item"><a class="slender-search-engine-link" href="javascript:void(0)" data-id="` + item.ID.String() + `" data-name="` + item.Name.HTMLString() + `" data-url="` + item.URL.HTMLString() + `" data-method="` + item.Method.HTMLString() + `" data-body="` + item.Body.HTMLString() + `">` + renderSearchEngineIcon(item.Name.String(), item.Icon.String()) + `<span class="slender-search-engine-text">` + item.Name.String() + `</span></a></li>`
	}

	return tpl
}

func renderSearchEngineIcon(name string, icon string) string {
	iconEle := ""

	i := icons.GetBuiltInIcon(icon)
	if i != "" {
		iconEle = `<img class="slender-search-engine-icon slender-built-in-icon" src="` + i + `" alt="icon" />`
	} else if icon != "" {
		iconEle = `<img class="slender-search-engine-icon" src="` + icon + `" alt="icon" />`
	} else {
		first := string([]rune(name)[0])
		iconEle = `<span class="slender-search-engine-icon slender-built-in-icon">` + first + `</span>`
	}

	return iconEle
}
