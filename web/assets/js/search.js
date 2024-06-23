(function () {
  'use strict';

  const inHomeSearchKey = 'slender-in-home-search';
  let inHomeSearch = localStorage.getItem(inHomeSearchKey, '0') === '1' ? true : false;
  const useSearchEngineKey = 'slender-use-search-engine';
  let useSearchEngine = localStorage.getItem(useSearchEngineKey, '0') || '0';

  const curSearchEngine = {
    id: '0',
    name: '',
    url: '',
    method: '',
    body: '',
    icon: '',
    builtInIcon: true,
  };

  const searchEngines = [];
  document.querySelectorAll('.slender-search-engine-item .slender-search-engine-link').forEach(ele => {
    let icon = '',
      builtInIcon = false;
    const img = ele.querySelector('.slender-search-engine-icon');
    if (img) {
      icon = img.src;
      if (img.classList.contains('slender-built-in-icon')) {
        builtInIcon = true;
      }
    }

    const data = ele.dataset;
    searchEngines.push({
      id: data.id,
      name: data.name,
      url: data.url,
      method: data.method,
      body: data.body,
      icon,
      builtInIcon,
    });
  });

  const hide = 'slender-search-result-hide';
  const showTip = 'slender-show-search-tip';
  const hideBtn = 'slender-hide-search-button';
  let timer = null;
  let debounced = null;

  function debounce(func, delay, immediate) {
    window.clearTimeout(timer);

    if (immediate) {
      if (timer === null) {
        func();
      }
    }

    timer = window.setTimeout(() => {
      func();
      timer = null;
    }, delay);

    return {
      cancel: () => {
        window.clearTimeout(timer);
        timer = null;
      },
    };
  }

  function restoreDisplay(bookmarkList) {
    if (debounced) {
      debounced.cancel();
      debounced = null;
    }

    bookmarkList.forEach(bm => {
      bm.ele.classList.remove(hide);
      bm.items.forEach(item => {
        item.ele.classList.remove(hide);
      });
    });
  }

  function renderSearchEngineIcon(name, icon, builtInIcon) {
    if (icon) {
      return `<img class="slender-search-input-prefix ${builtInIcon ? 'slender-built-in-icon' : ''}" src="${icon}" alt="icon" />`;
    } else {
      const first = name.length > 0 ? name[0] : '';
      return `<span class="slender-search-input-prefix slender-built-in-icon">${first}</span>`;
    }
  }

  function jumper(value, url, method, body) {
    if (method === 'GET') {
      window.open(url.replaceAll('%s', value));
    } else {
      let tempForm = document.createElement('form');
      tempForm.action = url;
      tempForm.target = '_blank';
      tempForm.method = 'post';
      tempForm.style.display = 'none';

      try {
        const params = JSON.parse(body);
        if (typeof params === 'object' && params) {
          for (const key in params) {
            const opt = document.createElement('textarea');
            opt.name = key;
            opt.value = params[key] === '%s' ? value : params[key];
            tempForm.appendChild(opt);
          }
        }

        document.body.appendChild(tempForm);
        tempForm.submit();
        tempForm.remove();
        tempForm = null;
      } catch (err) {
        console.error(err);
        tempForm.remove();
        tempForm = null;
      }
    }
  }

  const searchInput = document.getElementById('slender-search-input');
  const clearBtn = document.getElementById('slender-clear-button');
  const totalModule = document.getElementById('slender-search-total');
  let totalNum = 0;

  if (searchInput) {
    const curBookmarkList = [];
    document.querySelectorAll('.slender-folder-container').forEach(group => {
      const items = [];
      group.querySelectorAll('.slender-large-bookmark-item').forEach(item => {
        items.push({
          ele: item,
          url: item.querySelector('.slender-bookmark-link').href,
          title: item.querySelector('.slender-large-bookmark-title').innerText.trim().toLowerCase(),
          des: item.querySelector('.slender-large-bookmark-des').innerText.trim().toLowerCase(),
        });
        totalNum++;
      });

      group.querySelectorAll('.slender-bookmark-item').forEach(item => {
        items.push({
          ele: item,
          url: item.querySelector('.slender-bookmark-link').href,
          title: item.querySelector('.slender-bookmark-text').innerText.trim().toLowerCase(),
          des: item.querySelector('.slender-bookmark-link').title.trim().toLowerCase(),
        });
        totalNum++;
      });

      curBookmarkList.push({
        ele: group,
        items,
      });
    });

    function bookmarksHandler() {
      const val = searchInput.value.trim().toLowerCase();
      let curTotalNum = totalNum;

      curBookmarkList.forEach(bm => {
        let hadShowItem = false;

        bm.items.forEach(item => {
          if (!inHomeSearch && curSearchEngine.id != '0') {
            item.ele.classList.remove(hide);
          } else {
            if (item.url.includes(val) || item.title.includes(val) || item.des.includes(val)) {
              item.ele.classList.remove(hide);
              hadShowItem = true;
            } else {
              item.ele.classList.add(hide);
              curTotalNum--;
            }
          }
        });

        if (!inHomeSearch && curSearchEngine.id != '0') {
          bm.ele.classList.remove(hide);
        } else {
          if (hadShowItem) {
            bm.ele.classList.remove(hide);
          } else {
            bm.ele.classList.add(hide);
          }
        }
      });

      if (clearBtn) {
        if (val) {
          clearBtn.classList.add(showTip);
        } else {
          clearBtn.classList.remove(showTip);
        }
      }

      if (totalModule) {
        if (!inHomeSearch && curSearchEngine.id != '0') {
          totalModule.classList.remove(showTip);
        } else {
          if (val) {
            totalModule.textContent = `${curTotalNum}/${totalNum}`;
            totalModule.classList.add(showTip);
          } else {
            totalModule.classList.remove(showTip);
          }
        }
      }
    }

    const inHomeInput = document.getElementById('slender-search-enable-in-home');
    if (inHomeInput) {
      if (inHomeSearch) {
        inHomeInput.checked = true;
      }

      inHomeInput.addEventListener('change', e => {
        const checked = e.target.checked;
        inHomeSearch = checked;
        localStorage.setItem(inHomeSearchKey, inHomeSearch ? '1' : '0');
        bookmarksHandler();
      });
    }

    searchInput.addEventListener('keydown', ev => {
      if (ev.key === 'Escape') {
        restoreDisplay(curBookmarkList);
        searchInput.value = '';
        clearBtn && clearBtn.classList.remove(showTip);
        return false;
      } else if (ev.key === 'Enter' && curSearchEngine.id != '0' && searchInput.value) {
        jumper(searchInput.value, curSearchEngine.url, curSearchEngine.method, curSearchEngine.body);
        searchInput.value = '';
        clearBtn && clearBtn.classList.remove(showTip);
        return false;
      }

      debounced = debounce(bookmarksHandler, 250);
    });

    if (clearBtn) {
      clearBtn.addEventListener('click', ev => {
        ev.stopPropagation();
        restoreDisplay(curBookmarkList);
        searchInput.value = '';
        clearBtn.classList.remove(showTip);
        if (totalModule) {
          totalModule.classList.remove(showTip);
        }
      });
    }

    function searchIconHandler() {
      let findSearchEngine = false;
      for (const item of searchEngines) {
        if (useSearchEngine == item.id) {
          Object.assign(curSearchEngine, item);
          findSearchEngine = true;
          break;
        }
      }

      if (!findSearchEngine) {
        Object.assign(curSearchEngine, {
          id: '0',
          name: 'In Home Search',
          url: '',
          method: '',
          body: '',
          icon: '/assets/icons/mdi-homesearch.svg',
          builtInIcon: true,
        });
      }

      const iconEle = document.querySelector('.slender-search-input-prefix');
      if (iconEle) {
        iconEle.outerHTML = renderSearchEngineIcon(curSearchEngine.name, curSearchEngine.icon, curSearchEngine.builtInIcon);
      }
    }

    searchIconHandler();

    const showSearchEngineModule = 'slender-search-engine-module-show';
    const prefixBtn = document.getElementById('slender-search-input-prefix-overlay');
    const searchEngineModule = document.getElementById('slender-search-engine-module');
    if (prefixBtn && searchEngineModule) {
      prefixBtn.addEventListener('click', e => {
        e.stopPropagation();
        if (searchEngineModule.classList.contains(showSearchEngineModule)) {
          searchEngineModule.classList.remove(showSearchEngineModule);
          searchInput.focus();
          document.body.style.overflow = 'auto';
        } else {
          searchEngineModule.classList.add(showSearchEngineModule);
          document.body.style.overflow = 'hidden';
        }
      });

      searchEngineModule.addEventListener('click', e => {
        e.stopPropagation();
        if (e.target === e.currentTarget) {
          searchEngineModule.classList.remove(showSearchEngineModule);
          searchInput.focus();
          document.body.style.overflow = 'auto';
        }
      });
    }

    const searchIcon = document.querySelector('.slender-search-input-suffix');
    const searchBtn = document.getElementById('slender-search-input-suffix-overlay');

    if (searchBtn) {
      searchBtn.addEventListener('click', e => {
        e.stopPropagation();
        if (searchInput.value) {
          jumper(searchInput.value, curSearchEngine.url, curSearchEngine.method, curSearchEngine.body);
          searchInput.value = '';
          clearBtn && clearBtn.classList.remove(showTip);
        }
      });
    }

    function searchButtonHandler() {
      if (curSearchEngine.id == '0') {
        searchIcon && searchIcon.classList.add(hideBtn);
        searchBtn && searchBtn.classList.add(hideBtn);
        searchInput.classList.add(hideBtn);
      } else {
        searchIcon && searchIcon.classList.remove(hideBtn);
        searchBtn && searchBtn.classList.remove(hideBtn);
        searchInput.classList.remove(hideBtn);
      }
    }

    searchButtonHandler();

    const searchEngineList = document.getElementById('slender-search-engine-list');
    if (searchEngineList) {
      searchEngineList.addEventListener('click', e => {
        const target = e.target;
        if (target) {
          let searchEngineEle = null;

          if (target.classList.contains('slender-search-engine-link')) {
            searchEngineEle = target;
          } else {
            const parentEle = target.closest('.slender-search-engine-link');
            if (parentEle) {
              searchEngineEle = parentEle;
            }
          }

          if (searchEngineEle) {
            useSearchEngine = searchEngineEle.dataset.id;
            searchIconHandler();
            bookmarksHandler();
            searchButtonHandler();
            localStorage.setItem(useSearchEngineKey, useSearchEngine);
            searchEngineModule.classList.remove(showSearchEngineModule);
            searchInput.focus();
            document.body.style.overflow = 'auto';
          }
        }
      });
    }
  }
})();
