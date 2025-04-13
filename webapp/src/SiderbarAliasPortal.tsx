// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {useEffect, useState} from 'react';
import ReactDOM from 'react-dom';

interface Props {
    aliasFrom: string;
    aliasTo: string;
}

export default function SidebarAliases({aliasFrom, aliasTo}: Props) {
    const [target, setTarget] = useState<Element | null>(null);

    useEffect(() => {
        const item = document.querySelector(
            `a.SidebarLink[id^="sidebarItem_${aliasFrom}"] span.SidebarChannelLinkLabel`,
        );

        if (item) {
            setTarget(item);
        }
    }, [aliasFrom]);

    return target ? ReactDOM.createPortal(
        ` (${aliasTo})`,
        target,
    ) : null;
}
