package org.conf2.yang;

/**
 *
 */
public class LeafList extends MetaBase implements Describable {
    private String description;
    private String type;

    public LeafList(String ident) {
        super(ident);
    }

    public void setDescription(String description) {
        this.description = description;
    }

    @Override
    public String getDescription() {
        return description;
    }

    public void setType(String type) {
        this.type = type;
    }

    public String getType() {
        return type;
    }
}
