package org.conf2.schema;

/**
 *
 */
public abstract class RpcBase extends CollectionBase implements HasGroupings, HasTypedefs {
    private LinkedListCollection groupings = new LinkedListCollection("groupings", this);
    private LinkedListCollection typedefs = new LinkedListCollection("typedefs", this);

    public RpcBase(String ident) {
        super(ident);
    }

    @Override
    public MetaCollection getGroupings() {
        return groupings;
    }

    @Override
    public LinkedListCollection getTypedefs() {
        return typedefs;
    }

    @Override
    public void addMeta(Meta m) {
        if (m instanceof Grouping) {
            groupings.addMeta(m);
        } else if (m instanceof Typedef) {
            typedefs.addMeta(m);
        } else {
            super.addMeta(m);
        }
    }
}
