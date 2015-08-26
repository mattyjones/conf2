package org.conf2.yang;

/**
 *
 */
public class Grouping extends CollectionBase implements Describable, HasGroupings, HasTypedefs {
    private LinkedListCollection groupings = new LinkedListCollection("groupings", this);
    private LinkedListCollection typedefs = new LinkedListCollection("typedefs", this);
    public Grouping(String ident) {
        super(ident);
    }

    @Override
    public MetaCollection getGroupings() {
        return groupings;
    }

    @Override
    public MetaCollection getTypedefs() {
        return typedefs;
    }
}
