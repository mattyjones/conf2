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
}
