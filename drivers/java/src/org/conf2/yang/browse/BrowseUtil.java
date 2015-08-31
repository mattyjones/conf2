package org.conf2.yang.browse;

import org.conf2.yang.Meta;
import org.conf2.yang.MetaError;
import org.conf2.yang.driver.DriverError;

import java.lang.reflect.Field;
import java.lang.reflect.Method;

/**
 *
 */
public class BrowseUtil {

    public static final String accessorMethodNameFromMeta(String prefix, String ident) {
        return prefix + Character.toUpperCase(ident.charAt(0)) + ident.substring(1);
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

    public static final void setterField(Meta m, Object o, BrowseValue v) {
        String fieldName = m.getIdent();
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

    public static final void getterMethod(Meta m, Object o, BrowseValue v) {
        String methodName = accessorMethodNameFromMeta("get", m.getIdent());
        try {
            // TODO
            Method method = o.getClass().getMethod(methodName);
            Object result = method.invoke(o);
            v.str = result.toString();
        } catch (ReflectiveOperationException e) {
            throw new MetaError("Method not found", e);
        }
    }
}
