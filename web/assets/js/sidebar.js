(function () {
  'use strict';

  const hideSidebar = 'slender-hide-sidebar';
  const showSidebar = 'slender-show-sidebar';
  const highlightTitle = 'slender-highlight-title';

  const sidebarContainer = document.querySelector('.slender-sidebar-container');
  if (sidebarContainer) {
    let state = true;
    const doc = document.documentElement;
    if (doc.clientWidth <= 768) {
      state = false;
    }

    function isInViewport(ele) {
      const rect = ele.getBoundingClientRect();
      return (
        rect.top >= 0 &&
        rect.left >= 0 &&
        rect.bottom <= (window.innerHeight || document.documentElement.clientHeight) &&
        rect.right <= (window.innerWidth || document.documentElement.clientWidth)
      );
    }

    function focusAnimation(doc, ele) {
      const clientWidth = doc.clientWidth;
      if (clientWidth <= 768) {
        setSidebarState(false);
      }

      ele.classList.add(highlightTitle);
      const timer = window.setTimeout(() => {
        window.clearTimeout(timer);
        ele.classList.remove(highlightTitle);
      }, 1500);
    }

    function scrollHandler(doc, ele) {
      const inViewport = isInViewport(ele);
      const offsetTop = ele.offsetTop - 8;

      if (inViewport) {
        doc.scrollTo({
          top: offsetTop,
          behavior: 'smooth',
        });

        focusAnimation(doc, ele);
      } else {
        const observer = new IntersectionObserver(
          (entries, observer) => {
            entries.forEach(entry => {
              if (entry.isIntersecting) {
                observer.unobserve(entry.target);
                focusAnimation(doc, ele);
              }
            });
          },
          { threshold: 1.0 }
        );
        observer.observe(ele);

        doc.scrollTo({
          top: offsetTop,
          behavior: 'smooth',
        });
      }
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
          scrollHandler(doc, folderTitle);
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
