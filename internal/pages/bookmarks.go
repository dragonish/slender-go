package pages

import (
	"html/template"
	"slender/internal/database"
	"slender/internal/global"
	"slender/internal/icons"
	"slender/internal/model"
	"sort"
)

// generateBookmarks returns bookmarks and sidebar templates.
func generateBookmarks(dynamic *model.PageDynamicURL, privacy bool, ungrouped string, latest string, hot string) (template.HTML, template.HTML) {
	bookmarksTpl := ""
	sidebarTpl := ""

	folderList := make([]model.HomeFolderListItem, 0)
	fErr := database.GetHomeFolderList(&folderList)
	if fErr != nil {
		return template.HTML(bookmarksTpl), template.HTML(sidebarTpl)
	}

	inOtherNetwork := global.Config.InOtherNetwork(dynamic.Origin)

	bookmarkList := make([]model.HomeBookmarkListItem, 0)
	bErr := database.GetHomeBookmarkList(privacy, inOtherNetwork, &bookmarkList)
	if bErr != nil {
		return template.HTML(bookmarksTpl), template.HTML(sidebarTpl)
	}

	largeFolderList := make([]model.HomeFolderListItem, 0)
	generalFolderList := make([]model.HomeFolderListItem, 0)
	for _, folder := range folderList {
		if folder.Large {
			largeFolderList = append(largeFolderList, folder)
		} else {
			generalFolderList = append(generalFolderList, folder)
		}
	}

	bookmarks := make(map[model.MyInt64][]model.HomeBookmarkListItem)
	for _, bookmark := range bookmarkList {
		if bookmarks[bookmark.FolderID] == nil {
			bookmarks[bookmark.FolderID] = make([]model.HomeBookmarkListItem, 0)
			bookmarks[bookmark.FolderID] = append(bookmarks[bookmark.FolderID], bookmark)
		} else {
			bookmarks[bookmark.FolderID] = append(bookmarks[bookmark.FolderID], bookmark)
		}
	}

	for _, largeFolder := range largeFolderList {
		if len(bookmarks[largeFolder.ID]) > 0 {
			bookmarksTpl += renderBookmarkList(dynamic, &largeFolder, bookmarks[largeFolder.ID], inOtherNetwork)
			if global.Config.ShowSidebar {
				sidebarTpl += renderSidebar(&largeFolder)
			}
		}
	}

	if global.Config.ShowHot {
		hotBookmarkList := make([]model.HomeBookmarkListItem, 0)
		hErr := database.GetHomeHotBookmarkList(privacy, inOtherNetwork, global.Config.HotTotal, &hotBookmarkList)
		if hErr == nil && len(hotBookmarkList) > 0 {
			h := model.HomeFolderListItem{
				ID:    -2,
				Name:  model.MyString(hot),
				Des:   "",
				Large: false,
			}
			bookmarksTpl += renderBookmarkList(dynamic, &h, hotBookmarkList, inOtherNetwork)
			if global.Config.ShowSidebar {
				sidebarTpl += renderSidebar(&h)
			}
		}
	}

	if global.Config.ShowLatest {
		latestBookamrkList := make([]model.HomeBookmarkListItem, 0)
		lErr := database.GetHomeLatestBookmarkList(privacy, inOtherNetwork, global.Config.LatestTotal, &latestBookamrkList)
		if lErr == nil && len(latestBookamrkList) > 0 {
			l := model.HomeFolderListItem{
				ID:    -1,
				Name:  model.MyString(latest),
				Des:   "",
				Large: false,
			}
			bookmarksTpl += renderBookmarkList(dynamic, &l, latestBookamrkList, inOtherNetwork)
			if global.Config.ShowSidebar {
				sidebarTpl += renderSidebar(&l)
			}
		}
	}

	if len(bookmarks[0]) > 0 {
		u := model.HomeFolderListItem{
			ID:    0,
			Name:  model.MyString(ungrouped),
			Des:   "",
			Large: false,
		}
		bookmarksTpl += renderBookmarkList(dynamic, &u, bookmarks[0], inOtherNetwork)
		if global.Config.ShowSidebar {
			sidebarTpl += renderSidebar(&u)
		}
	}

	for _, generalFolder := range generalFolderList {
		if len(bookmarks[generalFolder.ID]) > 0 {
			bookmarksTpl += renderBookmarkList(dynamic, &generalFolder, bookmarks[generalFolder.ID], inOtherNetwork)
			if global.Config.ShowSidebar {
				sidebarTpl += renderSidebar(&generalFolder)
			}
		}
	}

	return template.HTML(bookmarksTpl), template.HTML(sidebarTpl)
}

func renderBookmarkList(dynamic *model.PageDynamicURL, folder *model.HomeFolderListItem, bookmarks []model.HomeBookmarkListItem, inOtherNetwork bool) string {
	tpl := ""

	switch folder.SortBy {
	case "visits":
		sort.SliceStable(bookmarks, func(i, j int) bool {
			return bookmarks[i].Visits > bookmarks[j].Visits
		})
	case "created_time":
		sort.SliceStable(bookmarks, func(i, j int) bool {
			return bookmarks[i].CreatedTime.After(bookmarks[j].CreatedTime)
		})
	}

	target := "_self"
	if global.Config.OpenInNewWindow {
		target = "_blank"
	}

	if folder.Large {
		for _, item := range bookmarks {
			useUrl := item.URL
			if !inOtherNetwork && item.Intranet != "" {
				useUrl = item.Intranet
			}
			url := dynamic.Convert(useUrl.String())
			name := item.Name.String()

			if name == "" {
				name = url
			}

			iconELe := ""
			icon := icons.GetBuiltInIcon(item.Icon.String())
			if icon != "" {
				iconELe = `<img class="slender-large-bookmark-icon slender-built-in-icon" src="` + icon + `" alt="icon" loading="lazy" />`
			} else if item.Icon != "" {
				iconELe = `<img class="slender-large-bookmark-icon" src="` + item.Icon.String() + `" alt="icon" loading="lazy" />`
			} else if global.Config.UseLetterIcon && name != "" {
				first := string([]rune(name)[0])
				iconELe = `<div class="slender-large-bookmark-icon slender-built-in-icon">` + first + `</div>`
			}

			des := item.Des.String()
			if des == "" {
				des = "&nbsp;"
			}

			tpl += `<div class="slender-large-bookmark-item"><a class="slender-bookmark-link slender-large-bookmark-link" href="` + url + `" target="` + target + `" rel="noopener" data-id="` + item.ID.String() + `">` + iconELe + `<div class="slender-large-bookmark-content"><span class="slender-large-bookmark-text slender-large-bookmark-title">` + name + `</span><span class="slender-large-bookmark-text slender-large-bookmark-des">` + des + `</span></div></a></div>`
		}

		tpl = `<div id="slender-folder-` + folder.ID.String() + `" class="slender-folder-container slender-large-folder"><h3 class="slender-bookmark-folder-title" data-folder="` + folder.ID.String() + `">` + folder.Name.String() + `</h3><div class="slender-large-bookmark-list">` + tpl + `</div></div>`
	} else {
		for _, item := range bookmarks {
			useUrl := item.URL
			if !inOtherNetwork && item.Intranet != "" {
				useUrl = item.Intranet
			}
			url := dynamic.Convert(useUrl.String())
			name := item.Name.String()

			if name == "" {
				name = url
			}

			iconELe := ""
			icon := icons.GetBuiltInIcon(item.Icon.String())
			if icon != "" {
				iconELe = `<img class="slender-bookmark-icon slender-built-in-icon" src="` + icon + `" alt="icon" loading="lazy" />`
			} else if item.Icon != "" {
				iconELe = `<img class="slender-bookmark-icon" src="` + item.Icon.String() + `" alt="icon" loading="lazy" />`
			} else if global.Config.UseLetterIcon && name != "" {
				first := string([]rune(name)[0])
				iconELe = `<span class="slender-bookmark-icon slender-built-in-icon">` + first + `</span>`
			}

			tpl += `<li class="slender-bookmark-item"><a class="slender-bookmark-link" href="` + url + `" target="` + target + `" rel="noopener" title="` + item.Des.String() + `" data-id="` + item.ID.String() + `">` + iconELe + `<span class="slender-bookmark-text">` + name + `</span></a></li>`
		}

		tpl = `<div id="slender-folder-` + folder.ID.String() + `" class="slender-folder-container"><h3 class="slender-bookmark-folder-title" data-folder="` + folder.ID.String() + `">` + folder.Name.String() + `</h3><ul class="slender-bookmark-list">` + tpl + `</ul></div>`
	}

	return tpl
}

func renderSidebar(folder *model.HomeFolderListItem) string {
	return `<li><a class="slender-sidebar-folder-item" href="#` + folder.ID.String() + `">` + folder.Name.String() + `</a></li>`
}
