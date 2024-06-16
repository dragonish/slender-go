(function () {
  'use strict';

  const hideSidebar = 'slender-hide-sidebar';
  const showSidebar = 'slender-show-sidebar';
  const sidebarContainer = document.querySelector('.slender-sidebar-container');
  if (sidebarContainer) {
    let state = true;
    const doc = document.documentElement;
    if (doc.clientWidth <= 768) {
      state = false;
    }

    function setSidebarState(sta) {
      if (sta !== undefined) {
        state = sta;
      } else {
        state = !state;
      }

      if (state) {
        sidebarContainer.classList.remove(hideSidebar);
        sidebarContainer.classList.add(showSidebar);
      } else {
        sidebarContainer.classList.remove(showSidebar);
        sidebarContainer.classList.add(hideSidebar);
      }
    }

    function blurSerachInput() {
      const searchInput = document.getElementById('slender-search-input');
      setTimeout(() => {
        searchInput?.blur();
      }, 0);
    }

    document.querySelectorAll('.slender-sidebar-container .slender-sidebar-list a').forEach(ele => {
      const id = ele.getAttribute('href').replace('#', '');
      ele.addEventListener('click', ev => {
        ev.stopPropagation();
        ev.preventDefault();
        const folderTitle = document.querySelector(`[data-folder="${id}"]`);
        if (folderTitle) {
          const offsetTop = folderTitle.offsetTop - 8;
          doc.scrollTo({
            top: offsetTop,
            behavior: 'smooth',
          });
          const clientWidth = doc.clientWidth;
          if (clientWidth <= 768) {
            setSidebarState(false);
          }
        } else {
          console.warn(`unable find folder title about id "${id}".`);
        }
      });
    });

    const sidebarBtn = document.getElementById('slender-sidebar-button');
    if (sidebarBtn) {
      sidebarBtn.addEventListener('click', ev => {
        ev.stopPropagation();
        ev.preventDefault();
        blurSerachInput();
        setSidebarState();
      });
      sidebarBtn.addEventListener(
        'touchstart',
        ev => {
          ev.stopPropagation();
          ev.preventDefault();
          blurSerachInput();
          setSidebarState();
        },
        {
          passive: false,
        }
      );
    }
  }
})();
