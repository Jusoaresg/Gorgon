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

    function updateEpisode(episodeId, tracking) {
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
            }

            if (target.classList.contains('cal-card-track')) {
                applyTracking(target, tracking);
                target.textContent = TRACKING_LABELS[tracking] || tracking || 'Unknown';
            }
        });
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
                updateEpisode(message.episodeID, message.tracking);
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
})();
