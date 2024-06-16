(function () {
  'use strict';

  const content = document.getElementById('slender-content');
  if (content) {
    content.addEventListener('click', e => {
      const target = e.target;
      if (target && target.classList.contains('slender-bookmark-link')) {
        const id = target.dataset.id;
        if (id) {
          fetch(`/api/v1/bookmarks/${id}/visits`, {
            method: 'POST',
          }).catch(err => {
            console.error(err);
          });
        }
      }
    });
  }
})();
