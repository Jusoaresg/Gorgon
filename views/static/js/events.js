(function () {
    'use strict';

    var WS_URL = (location.protocol === 'https:' ? 'wss://' : 'ws://') + location.host + '/api/v1/ws';

    var TRACKING_CLASSES = ['wanted', 'missing', 'skipped', 'snatched', 'downloaded'];

    var TRACKING_LABELS = {
        wanted: 'Wanted',
        missing: 'Missing',
        skipped: 'Skipped',
        snatched: 'Snatched',
        downloaded: 'Downloaded'
    };

    function removeTrackingClasses(el) {
        TRACKING_CLASSES.forEach(function (cls) {
            el.classList.remove(cls);
        });
    }

    function applyTracking(el, tracking) {
        removeTrackingClasses(el);
        if (TRACKING_CLASSES.indexOf(tracking) !== -1) {
            el.classList.add(tracking);
        }
    }

    function updateEpisode(episodeId, tracking, infoUrl) {
        var targets = document.querySelectorAll('[data-episode-id="' + episodeId + '"]');
        targets.forEach(function (target) {
            if (target.classList.contains('episode-card')) {
                applyTracking(target, tracking);

                var badge = target.querySelector('.status-badge');
                if (badge) {
                    applyTracking(badge, tracking);
                    badge.textContent = tracking;
                }

                var deleteBtn = target.querySelector('.btn-delete');
                if (deleteBtn) {
                    deleteBtn.style.display = tracking === 'downloaded' ? 'flex' : 'none';
                }

                var searchActions = target.querySelectorAll('.js-search-actions');
                searchActions.forEach(function (btn) {
                    btn.style.display = tracking === 'downloaded' ? 'none' : 'flex';
                });

                var shouldShowRelease = (tracking === 'downloaded' || tracking === 'snatched') && infoUrl;
                var releaseLink = target.querySelector('.js-release-link');
                if (shouldShowRelease) {
                    if (!releaseLink) {
                        var actionButtons = target.querySelector('.action-buttons');
                        if (actionButtons) {
                            releaseLink = document.createElement('a');
                            releaseLink.className = 'btn-icon-sm btn-release js-release-link';
                            releaseLink.target = '_blank';
                            releaseLink.title = 'View Release';
                            releaseLink.innerHTML = RELEASE_ICON;
                            actionButtons.appendChild(releaseLink);
                        }
                    }
                    if (releaseLink) {
                        releaseLink.href = infoUrl;
                        releaseLink.style.display = 'flex';
                    }
                } else if (releaseLink) {
                    releaseLink.style.display = 'none';
                }
            }

            if (target.classList.contains('cal-card-track')) {
                applyTracking(target, tracking);
                target.textContent = TRACKING_LABELS[tracking] || tracking || 'Unknown';
            }
        });
    }

    var TOAST_ICONS = {
        success: '<svg class="toast-icon" width="20" height="20" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>',
        error: '<svg class="toast-icon" width="20" height="20" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v4m0 4h.01M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z"></path></svg>',
        info: '<svg class="toast-icon" width="20" height="20" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>'
    };

    var CLOSE_ICON = '<svg class="toast-close" width="18" height="18" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>';

    var RELEASE_ICON = '<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"></path></svg>';

    function toastContainer() {
        var container = document.getElementById('toast-container');
        if (!container) {
            container = document.createElement('div');
            container.id = 'toast-container';
            document.body.appendChild(container);
        }
        container.classList.add('toast-container');
        return container;
    }

    function dismissToast(toast) {
        toast.classList.add('toast-removing');
        toast.addEventListener('animationend', function () {
            toast.remove();
        }, { once: true });
    }

    function displayToast(message, type) {
        var toast = document.createElement('div');
        toast.className = 'toast toast-' + (type || 'info');

        var icon = document.createElement('div');
        icon.innerHTML = TOAST_ICONS[type] || TOAST_ICONS.info;

        var text = document.createElement('span');
        text.className = 'toast-message';
        text.textContent = message;

        var close = document.createElement('button');
        close.type = 'button';
        close.className = 'toast-close';
        close.title = 'Dismiss';
        close.innerHTML = CLOSE_ICON;
        close.addEventListener('click', function () {
            dismissToast(toast);
        });

        toast.appendChild(icon);
        toast.appendChild(text);
        toast.appendChild(close);

        toastContainer().appendChild(toast);

        setTimeout(function () {
            dismissToast(toast);
        }, 4000);
    }

    function toNullProfileParams(event) {
        var params = event.detail.parameters;
        if (!params) return;
        if ('default_filter_profile_id' in params && params['default_filter_profile_id'] === '') {
            params['default_filter_profile_id'] = null;
        }
        if ('filter_profile_id' in params && params['filter_profile_id'] === '') {
            params['filter_profile_id'] = null;
        }
    }

    document.body.addEventListener('htmx:configRequest', toNullProfileParams);

    function showToast(event) {
        var message;
        var type = 'info';

        try {
            var response = JSON.parse(event.detail.xhr.response);
            var data = response.data;
            if (Array.isArray(data)) {
                data = data[0];
            }
            message = (data && data.toastMessage) || response.message || null;
        } catch (e) {
            message = null;
        }

        if (!message) {
            var status = event.detail.xhr.status;
            type = status >= 200 && status < 300 ? 'success' : 'error';
            message = type === 'success'
                ? 'Request successful'
                : status >= 500
                    ? 'Something went wrong on the server'
                    : 'Something went wrong';
        } else {
            type = event.detail.xhr.status >= 200 && event.detail.xhr.status < 300 ? 'success' : 'error';
        }

        displayToast(message, type);
    }

    function padEpisodeNumber(n) {
        return n < 10 ? '0' + n : '' + n;
    }

    function onSearchFinished(message) {
        var label = 'S' + message.season + 'E' + padEpisodeNumber(message.number);
        if (message.name) {
            label += ' ' + message.name;
        }
        switch (message.result) {
            case 'snatched':
                displayToast(label + ' — Snatched', 'success');
                break;
            case 'noResults':
                displayToast(label + ' — No results found', 'info');
                break;
            case 'notAired':
                displayToast(label + ' — Not aired yet', 'info');
                break;
            case 'error':
                displayToast(label + ' — Search failed', 'error');
                break;
        }
    }

    function connect() {
        var socket;
        try {
            socket = new WebSocket(WS_URL);
        } catch (e) {
            return;
        }

        socket.onmessage = function (event) {
            var message;
            try {
                message = JSON.parse(event.data);
            } catch (e) {
                return;
            }

            if (message && message.type === 'EpisodeTrackingUpdated') {
                updateEpisode(message.episodeID, message.tracking, message.infoUrl);
            }

            if (message && message.type === 'EpisodeSearchFinished') {
                onSearchFinished(message);
            }
        };

        socket.onclose = function () {
            setTimeout(connect, 3000);
        };

        socket.onerror = function () {
            socket.close();
        };
    }

    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', connect);
    } else {
        connect();
    }

    window.showToast = showToast;
    window.toast = displayToast;
})();
