package org.conf2.schema.browse;

import org.conf2.schema.Choice;
import org.conf2.schema.Meta;
import org.conf2.schema.MetaCollection;
import org.conf2.schema.MetaList;

/**
 *
 */
public class MyNode implements Node {
    public MetaCollection meta;
    public Meta position;
    public BrowseEnter Enter;
    public BrowseRead Read;
    public BrowseIterate Iterate;
    public BrowseEdit Edit;
    public BrowseEvent Event;
    public BrowseChoose Choose;
    public boolean insideList;
    public boolean found;

    @Override
    public Node Select(Selection sel, MetaCollection meta, boolean create) {
        return null;
    }

    @Override
    public BrowseValue Read(Selection sel, Meta meta) {
        return null;
    }

    @Override
    public Node Next(Selection sel, MetaList meta, boolean create, BrowseValue[] key, boolean isFirst) {
        return null;
    }

    @Override
    public void Write(Selection sel, Meta meta, BrowseValue v) {

    }

    @Override
    public void Event(Selection sel, DataEvent e) {

    }

    @Override
    public Meta Choose(Selection sel, Choice choice) {
        return null;
    }
}
