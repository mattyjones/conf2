package org.conf2.schema.browse;

import org.conf2.schema.*;
import org.conf2.schema.yang.YangModule;

/**
 *
 */
public class ModuleBrowser implements Browser {
    private Module yang = YangModule.YANG;
    public Module module;

    public ModuleBrowser(Module module) {
        this.module = module;
    }

    public Module getModule() {
        return yang;
    }

    @Override
    public Selection getRootSelector() {
        Selection s = new Selection();
        s.meta = yang;
        s.Enter = () -> {
            s.found = module != null;
            return enterModule(module);
        };
        s.Edit = (EditOperation op, BrowseValue val) -> {
            switch (op) {
                case CREATE_CHILD:
                    module = new Module("unknown");
                    break;
                case UPDATE_VALUE:
                    BrowseUtil.setterMethod(s.position, module, val);
                    break;
            }
        };
        return s;
    }

    Selection enterModule(final Module module) {
        Selection s = new Selection();
        s.Read = () -> BrowseUtil.getterMethod(s.position, module);
        s.Enter = () -> {
            String ident = s.position.getIdent();
            if (ident.equals("revision")) {
                s.found = (module.revision != null);
                if (s.found) {
                    return enterRevision(module.revision);
                }
            } else if (ident.equals("rpcs")) {
                return enterRpcCollection(module.getRpcs());
            } else if (ident.equals("notifications")) {
                return enterNotifications(module.getNotifications());
            } else if (ident.equals("groupings")) {
                return enterGroupings(module.getGroupings());
            } else if (ident.equals("typedefs")) {
                return enterTypedefs(module.getTypedefs());
            } else if (ident.equals("definitions")) {
                return enterDefinitions(module);
            }
            return null;
        };
        s.Edit = (EditOperation op, BrowseValue val) -> {
            String ident = s.position.getIdent();
            switch (op) {
                case CREATE_LIST:
                    // nothing to do
                    break;
                case CREATE_CHILD:
                    if (ident.equals("revision")) {
                        module.revision = new Revision("unknown");
                    }
                    break;
                case UPDATE_VALUE:
                    BrowseUtil.setterMethod(s.position, module, val);
                    break;
            }
        };
        return s;
    }

    Selection enterDefinitions(MetaCollection defs) {
        final Meta[] def = new Meta[1];
        Selection s = new Selection();
        s.Iterate = (BrowseValue[] keys, boolean isFirst) -> {
            return false;
        };
        s.Enter = () -> {
            if (def[0] == null) {
                s.found = false;
                //throw new DriverError("Definition not created yet");
            } else if (def[0] instanceof Leaf) {
                return enterLeaf((Leaf) def[0]);
            } else if (def[0] instanceof LeafList) {
                return enterLeafList((LeafList) def[0]);
            } else if (def[0] instanceof MetaList) {
                return enterList((MetaList) def[0]);
            } else if (def[0] instanceof Container) {
                return enterContainer((Container) def[0]);
            } else if (def[0] instanceof Uses) {
                return enterUses((Uses) def[0]);
            } else if (def[0] instanceof Choice) {
                return enterChoice((Choice) def[0]);
            }
            return null;
        };
        s.Edit = (EditOperation op, BrowseValue val) -> {
            switch (op) {
                case CREATE_CHILD:
                    String ident = s.position.getIdent();
                    if ("leaf".equals(ident)) {
                        def[0] = new Leaf("unknown");
                    } else if ("leaf-list".equals(ident)) {
                        def[0] = new LeafList("unknown");
                    } else if ("container".equals(ident)) {
                        def[0] = new Container("unknown");
                    } else if ("list".equals(ident)) {
                        def[0] = new MetaList("unknown");
                    } else if ("uses".equals(ident)) {
                        def[0] = new Uses("unknown");
                    } else if ("choice".equals(ident)) {
                        def[0] = new Choice("unknown");
                    }
                    break;
                case POST_CREATE_CHILD:
                    defs.addMeta(def[0]);
                    break;
                case UPDATE_VALUE:
                    BrowseUtil.setterMethod(s.position, def[0], val);
                    break;
            }
        };
        return s;
    }

    Selection enterLeaf(Leaf leaf) {
        Selection s = new Selection();
        s.Enter = () -> {
            String ident = s.position.getIdent();
            if ("type".equals(ident)) {
                s.found = (leaf.getDataType() != null);
                if (s.found) {
                    return enterDataType(leaf.getDataType());
                }
            }
            return null;
        };
        s.Edit = (EditOperation op, BrowseValue val) -> {
            String ident = s.position.getIdent();
            switch (op) {
                case CREATE_CHILD:
                    if ("type".equals(ident)) {
                        leaf.setDataType(new DataType("unknown"));
                    }
                    break;
                case UPDATE_VALUE:
                    BrowseUtil.setterMethod(s.position, leaf, val);
                    break;
            }
        };

        return s;
    }

    Selection enterDataType(DataType type) {
        Selection s = new Selection();
        s.Edit = (EditOperation op, BrowseValue val) -> {
            switch (op) {
                case UPDATE_VALUE:
                    if ("ident".equals(s.position.getIdent())) {
                        type.setIdent(val.str);
                    } else {
                        BrowseUtil.setterField(s.position, type, val);
                    }
                    break;
            }
        };

        return s;
    }

    Selection enterLeafList(LeafList leafList) {
        Selection s = new Selection();
        s.Enter = () -> {
            String ident = s.position.getIdent();
            if ("type".equals(ident)) {
                s.found = (leafList.getDataType() != null);
                if (s.found) {
                    return enterDataType(leafList.getDataType());
                }
            }
            return null;
        };
        s.Edit = (EditOperation op, BrowseValue val) -> {
            String ident = s.position.getIdent();
            switch (op) {
                case CREATE_CHILD:
                    if ("type".equals(ident)) {
                        leafList.setDataType(new DataType("unknown"));
                    }
                    break;
                case UPDATE_VALUE:
                    BrowseUtil.setterMethod(s.position, leafList, val);
                    break;
            }
        };

        return s;
    }

    Selection enterUses(Uses uses) {
        Selection s = new Selection();
        s.Enter = () -> {
            return null;
        };
        s.Edit = (EditOperation op, BrowseValue val) -> {
            switch (op) {
                case UPDATE_VALUE:
                    BrowseUtil.setterMethod(s.position, uses, val);
                    break;
            }
        };

        return s;
    }

    Selection enterChoice(Choice choice) {
        Selection s = new Selection();
        s.Enter = () -> {
            String ident = s.position.getIdent();
            if ("cases".equals(ident)) {
                return enterChoiceCases(choice);
            }

            return null;
        };
        s.Edit = (EditOperation op, BrowseValue val) -> {
            String ident = s.position.getIdent();
            switch (op) {
                case UPDATE_VALUE:
                    BrowseUtil.setterMethod(s.position, choice, val);
                    break;
            }
        };

        return s;
    }

    Selection enterChoiceCases(MetaCollection cases) {
        final ChoiceCase[] choiceCase = new ChoiceCase[1];
        Selection s = new Selection();
        s.Iterate = (BrowseValue[] key, boolean isFirst) -> {
            return false;
        };
        s.Enter = () -> {
            String ident = s.position.getIdent();
            if ("definitions".equals(ident)) {
                return enterDefinitions(choiceCase[0]);
            }
            return null;
        };
        s.Edit = (EditOperation op, BrowseValue val) -> {
            switch (op) {
                case CREATE_LIST_ITEM:
                    choiceCase[0] = new ChoiceCase("unknown");
                    break;
                case POST_CREATE_LIST_ITEM:
                    cases.addMeta(choiceCase[0]);
                    break;
                case UPDATE_VALUE:
                    BrowseUtil.setterMethod(s.position, choiceCase[0], val);
                    break;
            }
        };

        return s;
    }

    Selection enterNotifications(MetaCollection notifications) {
        final Notification[] notification = new Notification[1];
        Selection s = new Selection();
        s.Iterate = (BrowseValue[] key, boolean isFirst) -> {
            return false;
        };
        s.Enter = () -> {
            String ident = s.position.getIdent();
            if (ident.equals("groupings")) {
                return enterGroupings(notification[0].getGroupings());
            } else if (ident.equals("typedefs")) {
                return enterTypedefs(notification[0].getTypedefs());
            } else if (ident.equals("definitions")) {
                return enterDefinitions(notification[0]);
            }
            return null;
        };
        s.Edit = (EditOperation op, BrowseValue val) -> {
            switch (op) {
                case CREATE_LIST_ITEM:
                    notification[0] = new Notification("unknown");
                    break;
                case POST_CREATE_LIST_ITEM:
                    notifications.addMeta(notification[0]);
                    break;
                case CREATE_CHILD:
                    break;
                case UPDATE_VALUE:
                    BrowseUtil.setterMethod(s.position, notification[0], val);
                    break;
            }
        };

        return s;
    }

    Selection enterContainer(Container container) {
        Selection s = new Selection();
        s.Enter = () -> {
            String ident = s.position.getIdent();
            if (ident.equals("groupings")) {
                return enterGroupings(container.getGroupings());
            } else if (ident.equals("typedefs")) {
                return enterTypedefs(container.getTypedefs());
            } else if (ident.equals("definitions")) {
                return enterDefinitions(container);
            }
            return null;
        };
        s.Edit = (EditOperation op, BrowseValue val) -> {
            switch (op) {
                case UPDATE_VALUE:
                    BrowseUtil.setterMethod(s.position, container, val);
                    break;
            }
        };

        return s;
    }

    Selection enterGroupings(MetaCollection groupings) {
        Selection s = new Selection();
        final Grouping[] grouping = new Grouping[1];
        s.Iterate = (BrowseValue[] key, boolean isFirst) -> {
            return false;
        };
        s.Enter = () -> {
            String ident = s.position.getIdent();
            if (ident.equals("groupings")) {
                return enterGroupings(grouping[0].getGroupings());
            } else if (ident.equals("typedefs")) {
                return enterTypedefs(grouping[0].getTypedefs());
            } else if (ident.equals("definitions")) {
                return enterDefinitions(grouping[0]);
            }
            return null;
        };
        s.Edit = (EditOperation op, BrowseValue val) -> {
            switch (op) {
                case CREATE_LIST_ITEM:
                    grouping[0] = new Grouping("unknown");
                    break;
                case POST_CREATE_LIST_ITEM:
                    groupings.addMeta(grouping[0]);
                    break;
                case UPDATE_VALUE:
                    BrowseUtil.setterMethod(s.position, grouping[0], val);
                    break;
            }
        };

        return s;
    }

    Selection enterTypedefs(MetaCollection typedefs) {
        Selection s = new Selection();
        final Typedef[] typedef = new Typedef[1];
        s.Iterate = (BrowseValue[] key, boolean isFirst) -> {
            return false;
        };
        s.Enter = () -> {
            String ident = s.position.getIdent();
            if ("type".equals(ident)) {
                s.found  = (typedef[0].getDataType() != null);
                if (s.found) {
                    return enterDataType(typedef[0].getDataType());
                }
            }
            return null;
        };
        s.Edit = (EditOperation op, BrowseValue val) -> {
            switch (op) {
                case CREATE_LIST_ITEM:
                    typedef[0] = new Typedef("unknown");
                    break;
                case POST_CREATE_LIST_ITEM:
                    typedefs.addMeta(typedef[0]);
                    break;
                case CREATE_CHILD:
                    String ident = s.position.getIdent();
                    if ("type".equals(ident)) {
                        typedef[0].setDataType(new DataType("unknown"));
                    }
                    break;
                case UPDATE_VALUE:
                    BrowseUtil.setterMethod(s.position, typedef[0], val);
                    break;
            }
        };
        return s;
    }

    Selection enterList(MetaList list) {
        Selection s = new Selection();
        s.Enter = () -> {
            String ident = s.position.getIdent();
            if (ident.equals("groupings")) {
                return enterGroupings(list.getGroupings());
            } else if (ident.equals("typedefs")) {
                return enterTypedefs(list.getTypedefs());
            } else if (ident.equals("definitions")) {
                return enterDefinitions(list);
            }
            return null;
        };
        s.Edit = (EditOperation op, BrowseValue val) -> {
            switch (op) {
                case UPDATE_VALUE:
                    BrowseUtil.setterMethod(s.position, list, val);
                    break;
            }
        };

        return s;
    }

    Selection enterRpcCollection(MetaCollection rpcs) {
        final Rpc[] rpc = new Rpc[1];
        Selection s = new Selection();
        s.Iterate = (BrowseValue[] key, boolean isFirst) -> {
            return false;
        };
        s.Enter = () -> {
            String ident = s.position.getIdent();
            if ("input".equals(ident)) {
                s.found = (rpc[0].input != null);
                if (s.found) {
                    return enterRpcBase(rpc[0].input);
                }
            } else if ("output".equals(ident)) {
                s.found = (rpc[0].output != null);
                if (s.found) {
                    return enterRpcBase(rpc[0].output);
                }
            }
            return null;
        };
        s.Edit = (EditOperation op, BrowseValue val) -> {
            switch (op) {
                case CREATE_LIST_ITEM:
                    rpc[0] = new Rpc("unknown");
                    break;
                case POST_CREATE_LIST_ITEM:
                    rpcs.addMeta(rpc[0]);
                    break;
                case CREATE_CHILD:
                    String ident = s.position.getIdent();
                    if ("input".equals(ident)) {
                        rpc[0].input = new RpcInput();
                    } else if ("output".equals(ident)) {
                        rpc[0].output = new RpcOutput();
                    }
                    break;
                case POST_CREATE_CHILD:
                    break;
                case UPDATE_VALUE:
                    BrowseUtil.setterMethod(s.position, rpc[0], val);
                    break;
            }
        };
        return s;
    }

    Selection enterRpcBase(RpcBase rpcBase) {
        Selection s = new Selection();
        s.Enter = () -> {
            String ident = s.position.getIdent();
            if (ident.equals("groupings")) {
                return enterGroupings(rpcBase.getGroupings());
            } else if (ident.equals("typedefs")) {
                return enterTypedefs(rpcBase.getTypedefs());
            } else if (ident.equals("definitions")) {
                return enterDefinitions(rpcBase);
            }
            return null;
        };
        s.Edit = (EditOperation op, BrowseValue val) -> {
System.out.println("ModuleBrowser.java: enterRpcBase op=" + op + " position=" + s.position.getIdent());
            switch (op) {
                case UPDATE_VALUE:
                    BrowseUtil.setterMethod(s.position, rpcBase, val);
                    break;
            }
        };
        return s;
    }

    Selection enterRevision(Revision rev) {
        Selection s = new Selection();
        s.Read = () -> BrowseUtil.getterMethod(s.position, module);
        s.Edit = (EditOperation op, BrowseValue val) -> {
            String ident = s.position.getIdent();
            if (ident.equals("rev-date")) {
                rev.setIdent(val.str);
            } else {
                BrowseUtil.setterMethod(s.position, module, val);
            }
        };
        return s;
    }
}
