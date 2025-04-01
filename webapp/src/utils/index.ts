// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

export function injectSidebarAliases(aliases: Record<string, string>) {
    const tryInject = () => {
        const items = document.querySelectorAll('[data-testid^="dm-channel-"]');
        if (items.length === 0) {
            return;
        }

        items.forEach((item) => {
            const usernameEl = item.querySelector('span');
            if (!usernameEl) {
                return;
            }

            const username = usernameEl.textContent?.trim();
            if (!username) {
                return;
            }

            const allUsers = (window as any).store?.getState()?.entities?.users?.users || {};

            // eslint-disable-next-line guard-for-in
            for (const userId in aliases) {
                const user = allUsers[userId];
                if (user?.username === username && !usernameEl.textContent?.includes(aliases[userId])) {
                    usernameEl.textContent += ` (${aliases[userId]})`;
                }
            }
        });
    };

    const observer = new MutationObserver(tryInject);
    const sidebar = document.querySelector('#SidebarContainer');
    if (sidebar) {
        observer.observe(sidebar, {childList: true, subtree: true});
        tryInject();
    }
}
