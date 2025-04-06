// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

export interface PluginRegistry {
    registerPostTypeComponent(typeName: string, component: React.ElementType);
    registerRootComponent(component: React.ElementType);
    registerWebSocketEventHandler(eventName: string, callback: () => unknown);

    // Add more if needed from https://developers.mattermost.com/extend/plugins/webapp/reference
}
