// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React, {useEffect, useState} from 'react';
import ReactDOM from 'react-dom';
import {useStore} from 'react-redux';
import type {Store} from 'redux';

import type {GlobalState} from '@mattermost/types/store';

import {getCurrentChannel} from 'mattermost-redux/selectors/entities/channels';

import manifest from '@/manifest';

type Aliases = Record<string, string>;

export default function AliasRoot(initial: {aliases: Aliases}) {
    const store = useStore<GlobalState>() as Store<GlobalState>;
    const [aliases, setAliases] = useState(initial.aliases);
    const [portalTarget, setPortalTarget] = useState<HTMLElement | null>(null);
    const [aliasText, setAliasText] = useState<string | null>(null);

    const refreshAliases = async () => {
        const response = await fetch(`/plugins/${manifest.id}/api/v1/aliases`);
        if (response.ok) {
            const updated = await response.json();
            setAliases(updated);
        }
    };

    useEffect(() => {
        const handler = (event: MessageEvent) => {
            if (event.data?.type === 'alias_update') {
                refreshAliases();
            }
        };

        window.addEventListener('message', handler);
        return () => window.removeEventListener('message', handler);
    }, []);

    useEffect(() => {
        const unsubscribe = store.subscribe(() => {
            const state = store.getState();

            // eslint-disable-next-line @typescript-eslint/ban-ts-comment
            // @ts-expect-error
            const currentChannel = getCurrentChannel(state);
            if (!currentChannel) {
                setAliasText(null);
                return;
            }

            const userId = currentChannel?.teammate_id;
            if (userId && aliases[userId]) {
                setAliasText(aliases[userId]);

                const target = document.querySelector('#channelHeaderTitle');
                if (target instanceof HTMLElement) {
                    setPortalTarget(target);
                }
            } else {
                setAliasText(null);
            }
        });

        const observer = new MutationObserver(() => {
            const items = document.querySelectorAll('a.SidebarLink[id^="sidebarItem_"]');

            items.forEach((item) => {
                const id = item.getAttribute('id');
                const match = id?.match(/^sidebarItem_([^_]+)__/);
                const targetUserId = match?.[1];

                if (!targetUserId || !(targetUserId in aliases)) {
                    return;
                }

                const label = item.querySelector('span.SidebarChannelLinkLabel');
                const alias = aliases[targetUserId];

                if (label && !label.textContent?.includes(` (${alias})`)) {
                    label.textContent += ` (${alias})`;
                }
            });
        });

        observer.observe(document, {childList: true, subtree: true});

        return () => {
            unsubscribe();
            observer.disconnect();
        };
    }, [aliases, store]);

    if (!portalTarget || !aliasText) {
        return null;
    }

    return ReactDOM.createPortal(
        <span style={{color: '#aaa'}}>{`(${aliasText})`}</span>,
        portalTarget,
    );
}
