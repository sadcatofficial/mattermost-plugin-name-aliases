// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {injectSidebarAliases} from './utils';

import manifest from '@/manifest';

// Тип для карты псевдонимов
type Aliases = Record<string, string>;

export default class Plugin {
    private aliases: Aliases = {};

    public async initialize() {
        // Загружаем псевдонимы текущего пользователя с backend
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

        if (Object.keys(this.aliases).length > 0) {
            injectSidebarAliases(this.aliases);
        }
    }
}

declare global {
    interface Window {
        registerPlugin(pluginId: string, plugin: Plugin): void;
    }
}

window.registerPlugin(manifest.id, new Plugin());

