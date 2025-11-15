(function () {
  'use strict';

  const content = document.getElementById('slender-content');
  if (content) {
    content.addEventListener('click', e => {
      const target = e.target;
      if (target) {
        let linkEle = null;

        if (target.classList.contains('slender-bookmark-link')) {
          linkEle = target;
        } else {
          const parentEle = target.closest('.slender-bookmark-link');
          if (parentEle) {
            linkEle = parentEle;
          }
        }

        if (linkEle) {
          const id = linkEle.dataset.id;
          if (id) {
            fetch(`/api/v1/bookmarks/${id}/visits`, {
              method: 'POST',
            }).catch(err => {
              console.error(err);
            });
          }
        }
      }
    });
  }

  window.rainseeHomeEvent = function (top, bottom) {
    document.querySelector('body').style.padding = `${top}px 0 ${bottom}px 0`;
  };

  if (navigator.userAgent.includes('Firefox') && location.protocol === 'http:') {
    document.querySelectorAll('img[src*=".svg"][loading="lazy"]').forEach(img => {
      img.loading = 'eager';
    });
  }
})();
