// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';

import AliasRoot from './AliasRoot';
import type {PluginRegistry} from './types/mattermost-webapp';

import manifest from '@/manifest';

type Aliases = Record<string, string>;

export default class Plugin {
    private aliases: Aliases = {};

    public async initialize(registry: PluginRegistry) {
        try {
            const response = await fetch(`/plugins/${manifest.id}/api/v1/aliases`, {
                headers: {
                    'X-Requested-With': 'XMLHttpRequest',
                },
            });

            if (response.ok) {
                this.aliases = await response.json();
            }
        } catch (err) {
            // console.error('Alias plugin error:', err);
        }

        registry.registerRootComponent(() => <AliasRoot aliases={this.aliases}/>);
        registry.registerWebSocketEventHandler(`custom_${manifest.id}_alias_update`, () => {
            window.postMessage({type: 'alias_update'}, '*');
        });
    }
}

declare global {
    interface Window {
        registerPlugin(pluginId: string, plugin: Plugin): void;
    }
}

window.registerPlugin(manifest.id, new Plugin());

