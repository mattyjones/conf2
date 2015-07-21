package org.conf2.meta;

import java.util.Iterator;

/**
 *
 */
public class Uses extends MetaBase implements Describable, Proxy {
    private Grouping grouping;

    public Uses(String ident) {
        super(ident);
    }

    Grouping findGrouping(String ident) {
        if (grouping == null) {
            MetaCollection parent = getParent();
            while (parent != null) {
                if (parent instanceof HasGroupings) {
                    HasGroupings hg = (HasGroupings)parent;
                    Iterator<Meta> groupings = new MetaCollectionIterator(hg.getGroupings());
                    grouping = (Grouping) MetaUtil.findByIdent(groupings, getIdent());
                    if (grouping != null) {
                        break;
                    }
                }
                parent = parent.getParent();
            }
        }
        return grouping;
    }

    @Override
    public Iterator<Meta> resolveProxy() {
        Grouping g = findGrouping(getIdent());
        return g == null ? null : new MetaCollectionIterator(grouping);
    }
}
