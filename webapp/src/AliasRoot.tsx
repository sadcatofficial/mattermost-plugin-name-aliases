// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React, {useEffect, useState} from 'react';

import ChannelHeaderAliasPortal from './ChannelHeaderAliasProtal';
import SidebarAliasesPortal from './SiderbarAliasPortal';

import manifest from '@/manifest';

type Aliases = Record<string, string>;

export default function AliasRoot(initial: { aliases: Aliases }) {
    const [aliases, setAliases] = useState(initial.aliases);

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

    return (
        <>
            <ChannelHeaderAliasPortal aliases={aliases}/>
            {Object.entries(aliases).map((alias) => {
                return (
                    <SidebarAliasesPortal
                        key={alias[0]}
                        aliasFrom={alias[0]}
                        aliasTo={alias[1]}
                    />
                );
            })}
        </>
    );
}
