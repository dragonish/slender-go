(function () {
  'use strict';

  const hide = 'slender-search-result-hide';
  const showBtn = 'slender-show-search-button';
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

    searchInput.addEventListener('keydown', ev => {
      if (ev.key === 'Escape') {
        restoreDisplay(curBookmarkList);
        searchInput.value = '';
        clearBtn && clearBtn.classList.remove(showBtn);
        return false;
      }

      debounced = debounce(() => {
        const val = ev.target.value.trim().toLowerCase();
        let curTotalNum = totalNum;
        curBookmarkList.forEach(bm => {
          let hadShowItem = false;
          bm.items.forEach(item => {
            if (item.url.includes(val) || item.title.includes(val) || item.des.includes(val)) {
              item.ele.classList.remove(hide);
              hadShowItem = true;
            } else {
              item.ele.classList.add(hide);
              curTotalNum--;
            }
          });

          if (hadShowItem) {
            bm.ele.classList.remove(hide);
          } else {
            bm.ele.classList.add(hide);
          }
        });

        if (clearBtn) {
          if (val) {
            clearBtn.classList.add(showBtn);
          } else {
            clearBtn.classList.remove(showBtn);
          }
        }

        if (totalModule) {
          if (val) {
            totalModule.textContent = `${curTotalNum}/${totalNum}`;
            totalModule.classList.add(showBtn);
          } else {
            totalModule.classList.remove(showBtn);
          }
        }
      }, 250);
    });

    if (clearBtn) {
      clearBtn.addEventListener('click', ev => {
        ev.stopPropagation();
        restoreDisplay(curBookmarkList);
        searchInput.value = '';
        clearBtn.classList.remove(showBtn);
        if (totalModule) {
          totalModule.classList.remove(showBtn);
        }
      });
    }
  }
})();
