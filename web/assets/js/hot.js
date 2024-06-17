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
})();
