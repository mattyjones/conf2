package org.conf2.yang.browse;

import org.conf2.yang.*;
import org.conf2.yang.driver.DriverError;

import java.lang.reflect.Field;
import java.lang.reflect.Method;
import java.util.Arrays;
import java.util.Collection;
import java.util.Iterator;

/**
 *
 */
public class BrowseUtil {

    public static final String accessorMethodNameFromMeta(String prefix, String ident) {
        StringBuilder sb = new StringBuilder(prefix.length() + ident.length());
        boolean upper = false;
        if (prefix.length() > 0) {
            upper = true;
            sb.append(prefix);
        }
        for (int i = 0; i < ident.length(); i++) {
            char c = ident.charAt(i);
            if (c == '-') {
                upper = true;
            } else {
                if (upper) {
                    sb.append(Character.toUpperCase(c));
                } else {
                    sb.append(c);
                }
                upper = false;
            }
        }
        return sb.toString();
    }

    public static final void setterMethod(Meta m, Object o, BrowseValue v) {
        String methodName = accessorMethodNameFromMeta("set", m.getIdent());
        try {
            if (v.isList) {
                switch (v.valType) {
                    case INT32: {
                        Method method = o.getClass().getMethod(methodName, int[].class);
                        method.invoke(o, new Object[] {v.int32list});
                        break;
                    }
                    case STRING: {
                        Method method = o.getClass().getMethod(methodName, String[].class);
                        method.invoke(o, new Object[] {v.strlist});
                        break;
                    }
                    case BOOLEAN: {
                        Method method = o.getClass().getMethod(methodName, boolean[].class);
                        method.invoke(o, new Object[] {v.boollist});
                        break;
                    }
                    case EMPTY:
                        break;
                    default:
                        throw new DriverError("Format " + v.valType + " not supported");
                }
            } else {
                switch (v.valType) {
                    case INT32: {
                        Method method = o.getClass().getMethod(methodName, int.class);
                        method.invoke(o, v.int32);
                        break;
                    }
                    case STRING: {
                        Method method = o.getClass().getMethod(methodName, String.class);
                        method.invoke(o, v.str);
                        break;
                    }
                    case BOOLEAN: {
                        Method method = o.getClass().getMethod(methodName, boolean.class);
                        method.invoke(o, v.bool);
                        break;
                    }
                    case EMPTY:
                        break;
                    default:
                        throw new DriverError("Format " + v.valType + " not supported");
                }
            }
        } catch (ReflectiveOperationException e) {
            String msg = String.format("Method %s not found on class %s", methodName, o.getClass().getSimpleName());
            throw new MetaError(msg, e);
        }
    }

    public static final void readField(Meta m, Object o, BrowseValue v) {
        String fieldName = accessorMethodNameFromMeta("", m.getIdent());
        try {
            DataType t = ((HasDataType) m).getDataType();
            v.isList = (m instanceof LeafList);
            v.valType = t.valType;
            Field field = o.getClass().getField(fieldName);
            switch (v.valType) {
                case INT32: {
                    if (v.isList) {
                        v.int32list = (int[]) field.get(o);
                    } else {
                        v.int32 = field.getInt(o);
                    }
                    break;
                }
                case STRING: {
                    if (v.isList) {
                        v.strlist = coerseStringList(field.get(o));
                    } else {
                        v.str = field.get(o).toString();
                    }
                    break;
                }
                case BOOLEAN: {
                    if (v.isList) {
                        v.boollist = (boolean[]) field.get(o);
                    } else {
                        v.bool = field.getBoolean(o);
                    }
                    break;
                }
                case EMPTY:
                    break;
                default:
                    throw new DriverError("Format " + v.valType + " not supported");
            }
        } catch (ReflectiveOperationException e) {
            String msg = String.format("Field %s not found on class %s", fieldName, o.getClass().getSimpleName());
            throw new MetaError(msg, e);
        }
    }

    public static final void setterField(Meta m, Object o, BrowseValue v) {
        String fieldName = accessorMethodNameFromMeta("", m.getIdent());
        try {
            switch (v.valType) {
                case INT32: {
                    Field field = o.getClass().getField(fieldName);
                    if (v.isList) {
                        field.set(o, v.int32list);
                    } else {
                        field.setInt(o, v.int32);
                    }
                    break;
                }
                case STRING: {
                    Field field = o.getClass().getField(fieldName);
                    if (v.isList) {
                        field.set(o, v.strlist);
                    } else {
                        field.set(o, v.str);
                    }
                    break;
                }
                case BOOLEAN: {
                    Field field = o.getClass().getField(fieldName);
                    if (v.isList) {
                        field.set(o, v.boollist);
                    } else {
                        field.setBoolean(o, v.bool);
                    }
                    break;
                }
                case EMPTY:
                    break;
                default:
                    throw new DriverError("Format " + v.valType + " not supported");
            }
        } catch (ReflectiveOperationException e) {
            String msg = String.format("Field %s not found on class %s", fieldName, o.getClass().getSimpleName());
            throw new MetaError(msg, e);
        }
    }

    public static String[] coerseStringList(Object o) {
        if (o instanceof Collection) {
            Collection c = (Collection) o;
            String[] strlist = new String[c.size()];
            Iterator items = c.iterator();
            for (int i = 0; items.hasNext(); i++) {
                Object item = items.next();
                strlist[i] = (item != null ? item.toString() : null);
            }
            return strlist;
        }

        return (String[])o;
    }

    public static final void getterMethod(Meta m, Object o, BrowseValue v) {
        String methodName = accessorMethodNameFromMeta("get", m.getIdent());
        try {
            DataType t = ((HasDataType) m).getDataType();
            v.isList = (m instanceof LeafList);
            Method method = o.getClass().getMethod(methodName);
            Object result = method.invoke(o);
            if (result == null) {
                return;
            }
            v.valType = t.valType;
            switch (v.valType) {
                case INT32: {
                    if (v.isList) {
                        v.int32list = (int[]) result;
                    } else {
                        v.int32 = (Integer) result;
                    }
                    break;
                }
                case STRING: {
                    if (v.isList) {
                        v.strlist = coerseStringList(result);
                    } else {
                        v.str = result.toString();
                    }
                    break;
                }
                case BOOLEAN: {
                    if (v.isList) {
                        v.boollist = (boolean[]) result;
                    } else {
                        v.bool = (Boolean) result;
                    }
                    break;
                }
                case EMPTY:
                    break;
                default:
                    throw new DriverError("Format " + v.valType + " not supported");
            }
        } catch (ReflectiveOperationException e) {
            throw new MetaError("Method not found", e);
        }
    }

    public static final String[] strArrayAppend(String[] strlist, String s) {
        if (strlist == null) {
            return new String[] { s };
        }
        String[] strlistNew = Arrays.copyOf(strlist, strlist.length + 1);
        strlistNew[strlist.length - 1] = s;
        return strlist;
    }

    public static final int[] intArrayAppend(int[] intlist, int n) {
        if (intlist == null) {
            return new int[] { n };
        }
        int[] intlistNew = Arrays.copyOf(intlist, intlist.length + 1);
        intlistNew[intlistNew.length - 1] = n;
        return intlistNew;
    }
}
