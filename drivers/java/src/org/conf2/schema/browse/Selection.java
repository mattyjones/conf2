package org.conf2.schema.browse;

import org.conf2.schema.MetaCollection;

/**
 *
 */
public class Selection {
    public Selection(Node node, MetaCollection meta) {
        state = new WalkState();
        state.node = node;
        state.meta = meta;
    }
    public WalkState state;
    public Node node;

    public Selection selectContainer(Node node) {
        return new Selection(node, (MetaCollection) state.position);
    }

    public Selection selectListItem(Node node, BrowseValue[] key) {
        Selection next = new Selection(node, (MetaCollection) state.position);
        next.state.insideList = true;
        if (key != null && key.length > 0) {
            next.state.key = key;
            //next.state.path = next.state.path.SetKey(key);
        }
        return next;
    }
}
