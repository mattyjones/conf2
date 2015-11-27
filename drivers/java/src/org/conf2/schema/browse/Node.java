package org.conf2.schema.browse;

import org.conf2.schema.Choice;
import org.conf2.schema.Meta;
import org.conf2.schema.MetaCollection;
import org.conf2.schema.MetaList;

/**
 *
 */
public interface Node {
    Node Select(Selection sel, MetaCollection meta, boolean create);
    BrowseValue Read(Selection sel, Meta meta);
    Node Next(Selection sel, MetaList meta, boolean create, BrowseValue[] key, boolean isFirst);
    void Write(Selection sel, Meta meta, BrowseValue v);
    void Event(Selection sel, DataEvent e);
    Meta Choose(Selection sel, Choice choice);
}
