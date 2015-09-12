package org.conf2.schema.browse;

import org.conf2.schema.Meta;
import org.conf2.schema.MetaCollection;
import org.conf2.schema.driver.DriverHandle;

/**
 *
 */
public class Selection {
    public MetaCollection meta;
    public Meta position;
    public BrowseEnter Enter;
    public BrowseRead Read;
    public BrowseIterate Iterate;
    public BrowseEdit Edit;
    public BrowseExit Exit;
    public BrowseChoose Choose;
    public boolean insideList;
    public boolean found;
}
