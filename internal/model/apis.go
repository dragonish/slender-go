package model

// ListCondition defines base page condition.
type ListCondition struct {
	Page  int    `json:"page"`  // number of pages.
	Size  int    `json:"size"`  // size per page.
	Order string `json:"order"` // sort condition.
}

// PageData defines page data.
type PageData struct {
	PageSize int     `json:"pageSize"` // size per page.
	PageNo   int     `json:"pageNo"`   // current page number.
	Total    MyInt64 `json:"total"`
}

// BookmarkPostBody defines add bookmark body.
type BookmarkPostBody struct {
	URL      MyString  `json:"url" db:"url"`
	Name     MyString  `json:"name" db:"name"`
	Des      MyString  `json:"description" db:"description"`
	Icon     MyString  `json:"icon" db:"icon"`
	Privacy  MyBool    `json:"privacy" db:"privacy"`
	Weight   MyInt16   `json:"weight" db:"weight"`
	FolderID NullInt64 `json:"folderId" db:"folder_id"`

	HideInOther MyBool `json:"hideInOther" db:"hide_in_other"`

	Files []MyInt64 `json:"files"`
}

// BookmarkPatchBody defines update bookmark body.
type BookmarkPatchBody struct {
	URL      *MyString  `json:"url,omitempty" db:"url"`
	Name     *MyString  `json:"name,omitempty" db:"name"`
	Des      *MyString  `json:"description,omitempty" db:"description"`
	Icon     *MyString  `json:"icon,omitempty" db:"icon"`
	Privacy  *MyBool    `json:"privacy,omitempty" db:"privacy"`
	Weight   *MyInt16   `json:"weight,omitempty" db:"weight"`
	Visits   *MyInt64   `json:"visits,omitempty" db:"visits"`
	FolderID *NullInt64 `json:"folderId,omitempty" db:"folder_id"`

	HideInOther *MyBool `json:"hideInOther,omitempty" db:"hide_in_other"`

	Files []MyInt64 `json:"files,omitempty"`
}

// BookmarkListCondition defines bookmark list condition.
type BookmarkListCondition struct {
	ListCondition

	Name MyString `json:"name"`
	Des  MyString `json:"description"`
	URL  MyString `json:"url"`

	Privacy *MyBool    `json:"privacy,omitempty"`
	Folder  *NullInt64 `json:"folder,omitempty"`
}

// FolderListCondition defines folder list condition.
type FolderListCondition struct {
	ListCondition

	Name MyString `json:"name"`
	Des  MyString `json:"description"`

	Privacy *MyBool `json:"privacy,omitempty"`
}

// BookmarkBaseData defines bookmark base data.
type BookmarkBaseData struct {
	ID   MyInt64  `json:"id" db:"id"`
	URL  MyString `json:"url" db:"url"`
	Name MyString `json:"name" db:"name"`
	Des  MyString `json:"description" db:"description"`
	Icon MyString `json:"icon" db:"icon"`

	Privacy MyBool  `json:"privacy" db:"privacy"`
	Weight  MyInt16 `json:"weight" db:"weight"`

	CreatedTime  MyDatetimeString `json:"createdTime" db:"created_time"`
	ModifiedTime MyDatetimeString `json:"modifiedTime" db:"modified_time"`

	Visits MyInt64 `json:"visits" db:"visits"`

	FolderID NullInt64 `json:"folderId" db:"folder_id"` // folder id.

	HideInOther MyBool `json:"hideInOther" db:"hide_in_other"` // hide in other network environments.
}

// BookmarkListItem defines bookmark list item.
type BookmarkListItem struct {
	BookmarkBaseData

	FolderName MyString `json:"folderName" db:"folder_name"` // folder name.
}

// BookmarkListData defines bookmark list data.
type BookmarkListData struct {
	PageData

	List []BookmarkListItem `json:"list"` // bookmark list.
}

// BookmarkResData defines bookmark response data.
type BookmarkResData struct {
	BookmarkBaseData

	Files []FileBaseData `json:"files"` // file list.
}

// BookmarkFolderInfo defines bookmark's folder info.
type BookmarkFolderInfo struct {
	ID      MyInt64  `json:"id" db:"id"`
	Name    MyString `json:"name" db:"name"`
	Privacy MyBool   `json:"privacy" db:"privacy"`
}

// FileBaseData defines file base data.
type FileBaseData struct {
	ID   MyInt64  `json:"id" db:"id"`
	Path MyString `json:"path" db:"path"`
}

// FileInfo defines file info.
type FileInfo struct {
	Path       MyString  `db:"path"`
	BookmarkID NullInt64 `db:"bookmark_id"`
}

// FolderPostBody defines add folder body.
type FolderPostBody struct {
	Name    MyString `json:"name" db:"name"`
	Des     MyString `json:"description" db:"description"`
	Large   MyBool   `json:"large" db:"large"`
	Privacy MyBool   `json:"privacy" db:"privacy"`
	Weight  MyInt16  `json:"weight" db:"weight"`
}

// FolderPatchBody defines update folder body.
type FolderPatchBody struct {
	Name    *MyString `json:"name,omitempty" db:"name"`
	Des     *MyString `json:"description,omitempty" db:"description"`
	Large   *MyBool   `json:"large,omitempty" db:"large"`
	Privacy *MyBool   `json:"privacy,omitempty" db:"privacy"`
	Weight  *MyInt16  `json:"weight,omitempty" db:"weight"`
}

// FolderBaseData defines folder base data.
type FolderBaseData struct {
	ID   MyInt64  `json:"id" db:"id"`
	Name MyString `json:"name" db:"name"`
	Des  MyString `json:"description" db:"description"`

	Large   MyBool  `json:"large" db:"large"`
	Privacy MyBool  `json:"privacy" db:"privacy"`
	Weight  MyInt16 `json:"weight" db:"weight"`

	CreatedTime  MyDatetimeString `json:"createdTime" db:"created_time"`
	ModifiedTime MyDatetimeString `json:"modifiedTime" db:"modified_time"`
}

// FolderListItem defines folder list item.
type FolderListItem struct {
	FolderBaseData

	BookmarkTotal MyInt64 `json:"bookmarkTotal" db:"bookmark_total"` // total number of bookmarks included.
}

// FolderListData defines folder list data.
type FolderListData struct {
	PageData

	List []FolderListItem `json:"list"` // folder list.
}

// RequestAdminPostBoby defines request administrator post body.
type RequestAdminPostBoby struct {
	Password string `json:"password"`
}

// ConfigPatchBody defines config patch body.
type ConfigPatchBody struct {
	Title        *string `json:"title,omitempty"`        // website title.
	CustomFooter *string `json:"customFooter,omitempty"` // custom footer.

	ShowSidebar     *bool `json:"showSidebar,omitempty"`     // show folders sidebar.
	ShowSearchInput *bool `json:"showSearchInput,omitempty"` // show search input.
	ShowScrollTop   *bool `json:"showScrollTop,omitempty"`   // show scroll to top button.

	ShowLatest  *bool  `json:"showLatest,omitempty"`  // show the module of latest added bookmarks.
	LatestTotal *uint8 `json:"latestTotal,omitempty"` // number of bookmarks in the latest module.
	ShowHot     *bool  `json:"showHot,omitempty"`     // show the module of hot bookmarks.
	HotTotal    *uint8 `json:"hotTotal,omitempty"`    // number of bookmarks in the hot module.

	UseLetterIcon   *bool `json:"useLetterIcon,omitempty"`   // use first letter as icon.
	OpenInNewWindow *bool `json:"openInNewWindow,omitempty"` // always open the bookmark in the new window.
}

// BatchPatchBody defines batch patch body.
type BatchPatchBody struct {
	DataSet []int64     `json:"dataSet"`
	Action  string      `json:"action"`
	Payload interface{} `json:"payload"`
}

// BookmarkImportItem defines import bookmark item.
type BookmarkImportItem struct {
	URL     MyString `json:"url" db:"url"`
	Name    MyString `json:"name" db:"name"`
	Des     MyString `json:"description" db:"description"`
	Icon    MyString `json:"icon" db:"icon"`
	Privacy MyBool   `json:"privacy" db:"privacy"`
	Weight  MyInt16  `json:"weight" db:"weight"`

	HideInOther MyBool `json:"hideInOther"`
}

// LoginListCondition defines login list condition.
type LoginListCondition struct {
	ListCondition

	IP MyString `json:"ip"`
	UA MyString `json:"ua"`

	Admin  *MyBool `json:"admin,omitempty"`
	Active *MyBool `json:"active,omitempty"`
}

// LoginListItem defines login list item.
type LoginListItem struct {
	LoginID   MyString         `json:"loginId" db:"login_id"`
	LoginTime MyDatetimeString `json:"loginTime" db:"login_time"`
	IP        MyString         `json:"ip" db:"ip"`
	UA        MyString         `json:"ua" db:"ua"`
	IsAdmin   MyBool           `json:"isAdmin" db:"is_admin"`

	MaxAge MyUint16     `json:"maxAge" db:"max_age"`
	Active NullableBool `json:"active" db:"active"`
}

// LoginListData defines login list data.
type LoginListData struct {
	PageData

	List []LoginListItem `json:"list"`
}

// FileListCondition defines file list condition.
type FileListCondition struct {
	ListCondition

	Path MyString `json:"path"`

	Use *MyBool `json:"use"`
}

// FileListItem defines file list item.
type FileListItem struct {
	FileBaseData

	Used MyBool `json:"used" db:"used"`
}

// FileListData defines file list data.
type FileListData struct {
	PageData

	List []FileListItem `json:"list"`
}

// AboutInfoData defines program info.
type AboutInfoData struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	OS        string `json:"os"`
	Arch      string `json:"arch"`
	BuildDate string `json:"buildDate"`
}

// AboutIconsData defines icons info.
type AboutIconsData struct {
	MDIVersion string `json:"mdiVersion"`
	SIVersion  string `json:"siVersion"`
}

// SearchEnginePostBody defines add search engine body.
type SearchEnginePostBody struct {
	Name   MyString `json:"name" db:"name"`
	Method MyString `json:"method" db:"method"`
	URL    MyString `json:"url" db:"url"`
	Body   MyString `json:"body" db:"body"`
	Icon   MyString `json:"icon" db:"icon"`
	Weight MyInt16  `json:"weight" db:"weight"`
}

// SearchEngineListCondition defines search engine list condition.
type SearchEngineListCondition struct {
	ListCondition

	Name MyString `json:"name"`
	URL  MyString `json:"url"`

	Method *MyString `json:"method,omitempty"`
}

// SearchEngineBaseData defines search engine base data.
type SearchEngineBaseData struct {
	ID     MyInt64  `json:"id" db:"id"`
	Name   MyString `json:"name" db:"name"`
	Method MyString `json:"method" db:"method"`
	URL    MyString `json:"url" db:"url"`
	Body   MyString `json:"body" db:"body"`
	Icon   MyString `json:"icon" db:"icon"`
	Weight MyInt16  `json:"weight" db:"weight"`

	CreatedTime  MyDatetimeString `json:"createdTime" db:"created_time"`
	ModifiedTime MyDatetimeString `json:"modifiedTime" db:"modified_time"`
}

// SearchEngineListData defines search engine list data.
type SearchEngineListData struct {
	PageData

	List []SearchEngineBaseData `json:"list"`
}

// SearchEnginePatchBody defines update search engine body.
type SearchEnginePatchBody struct {
	Name   *MyString `json:"name,omitempty" db:"name"`
	Method *MyString `json:"method,omitempty" db:"method"`
	URL    *MyString `json:"url,omitempty" db:"url"`
	Body   *MyString `json:"body,omitempty" db:"body"`
	Icon   *MyString `json:"icon,omitempty" db:"icon"`
	Weight *MyInt16  `json:"weight,omitempty" db:"weight"`
}
