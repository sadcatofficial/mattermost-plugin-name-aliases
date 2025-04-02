// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

/* eslint-disable no-console */

export function injectSidebarAliases(aliases: Record<string, string>) {
    const tryInject = () => {
        const items = document.querySelectorAll('[id^="sidebarItem_"]');
        console.log({items});
        items.forEach((item) => {
            const id = item.getAttribute('id');
            if (!id) {
                return;
            }

            const match = id.split('_');
            const targetUserId = match?.[1];
            if (!targetUserId || !(targetUserId in aliases)) {
                return;
            }

            // const label = item.querySelector('.SidebarChannelLinkLabel');
            // if (!label || label.textContent?.includes(aliases[targetUserId])) {
            //     return;
            // }
            console.log('дописываем значение');
            item.textContent += ` (${aliases[targetUserId]})`;
        });
    };

    const observer = new MutationObserver(tryInject);
    const sidebar = document.querySelector('#SidebarContainer');
    if (sidebar) {
        observer.observe(sidebar, {childList: true, subtree: true});
        tryInject();
    }
}
