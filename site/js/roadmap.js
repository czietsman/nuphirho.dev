/* nuphirho.dev -- roadmap.js
 * Renders the publishing calendar from calendar-data.json.
 */
(function () {
  'use strict';

  var DAY_ORDER = ['mon','tue','wed','thu','fri','sat','sun'];
  var DAY_LABELS = {mon:'Mon',tue:'Tue',wed:'Wed',thu:'Thu',fri:'Fri',sat:'Sat',sun:'Sun'};
  var DAY_CLASS = {mon:'dark',tue:'dark',wed:'wed',thu:'dark',fri:'fri',sat:'weekend',sun:'weekend'};

  function buildPost(entry) {
    if (!entry) return '';
    var status = entry.status;
    var cat = entry.category;
    var catClass = status === 'published' ? 'published ' + cat :
                   status === 'pending-crosspost' ? 'pending-crosspost' :
                   'planned ' + cat;
    var check = (status === 'published') ? '<span class="check">\u2713</span>' : '';
    var titleHtml = entry.title
      ? '<span class="title">' + entry.title + '</span>'
      : '<span class="pending-title">Coming ' + (entry.label ? '' : 'soon') + '</span>';
    var labelHtml = '<span class="cat">' + (entry.label || entry.category) + '</span>';
    var platformsHtml = '';
    if (entry.platforms && entry.platforms.length) {
      var icons = entry.platforms.map(function (p) {
        if (p.pending || !p.url) return '<span class="platform-icon pending">' + p.label + '</span>';
        return '<a class="platform-icon" href="' + p.url + '" target="_blank">' + p.label + '</a>';
      }).join('');
      platformsHtml = '<div class="platforms">' + icons + '</div>';
    }
    return '<div class="post ' + catClass + '">' + check + labelHtml + titleHtml + platformsHtml + '</div>';
  }

  function buildWeek(week) {
    var days = week.days || {};
    var cols = '<div class="week-label"><div class="wk">Week</div><div class="dates">' + week.label + '</div></div>';
    for (var i = 0; i < DAY_ORDER.length; i++) {
      var d = DAY_ORDER[i];
      var cls = DAY_CLASS[d];
      var post = buildPost(days[d]);
      cols += '<div class="day ' + cls + '"><div class="day-name">' + DAY_LABELS[d] + '</div>' + post + '</div>';
    }
    return '<div class="week">' + cols + '</div>';
  }

  function buildMonth(month) {
    var isFuture = month.future;
    var collapsed = isFuture ? ' collapsed' : '';
    var toggleText = isFuture ? '\u25B8 Show' : '\u25BE Hide';
    var weeks = (month.weeks || []).map(buildWeek).join('');
    return '<div class="month-section">' +
      '<div class="month-header" onclick="nuphirhoToggleMonth(this)">' +
        '<span>' + month.label + '</span>' +
        '<span class="toggle">' + toggleText + '</span>' +
      '</div>' +
      '<div class="month-content' + collapsed + '">' + weeks + '</div>' +
    '</div>';
  }

  window.nuphirhoToggleMonth = function (header) {
    var content = header.nextElementSibling;
    var toggle = header.querySelector('.toggle');
    var isCollapsed = content.classList.contains('collapsed');
    content.classList.toggle('collapsed');
    toggle.textContent = isCollapsed ? '\u25BE Hide' : '\u25B8 Show';
  };

  fetch('/roadmap/calendar-data.json')
    .then(function (r) { return r.json(); })
    .then(function (data) {
      document.getElementById('calendar-root').innerHTML =
        data.months.map(buildMonth).join('');
    })
    .catch(function () {
      document.getElementById('calendar-root').innerHTML =
        '<p style="color:var(--roadmap-muted);font-family:\'DM Mono\',monospace;font-size:0.8rem;padding:1rem;">Calendar data unavailable.</p>';
    });
})();
