// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React, {useEffect, useState} from 'react';
import ReactDOM from 'react-dom';
import {useStore} from 'react-redux';
import type {Store} from 'redux';

import type {GlobalState} from '@mattermost/types/store';

import {getCurrentChannel} from 'mattermost-redux/selectors/entities/channels';

type Props = {
    aliases: Record<string, string>;
};

export default function ChannelHeaderAliasPortal({aliases}: Props) {
    const store = useStore<GlobalState>() as Store<GlobalState>;
    const [portalTarget, setPortalTarget] = useState<HTMLElement | null>(null);
    const [aliasText, setAliasText] = useState<string | null>(null);

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

        return () => unsubscribe();
    }, [aliases, store]);

    if (!portalTarget || !aliasText || !document.body.contains(portalTarget)) {
        return null;
    }

    return ReactDOM.createPortal(
        <span style={{color: '#aaa'}}>
            {`(${aliasText})`}
        </span>,
        portalTarget,
    );
}
