import React, { InputHTMLAttributes, forwardRef } from 'react';
import ErrorMessage from './ErrorMessage';

interface InputFieldProps extends InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  error?: string;
  helperText?: string;
  fullWidth?: boolean;
  icon?: React.ReactNode;
  iconPosition?: 'left' | 'right';
}

const InputField = forwardRef<HTMLInputElement, InputFieldProps>(
  (
    {
      label,
      error,
      helperText,
      fullWidth = true,
      className = '',
      icon,
      iconPosition = 'left',
      id,
      ...rest
    },
    ref
  ) => {
    const inputId = id || `input-${Math.random().toString(36).substr(2, 9)}`;
    const widthClass = fullWidth ? 'w-full' : '';
    const errorClass = error ? 'border-red-500 focus:ring-red-500 focus:border-red-500' : 'border-gray-300 focus:ring-primary-500 focus:border-primary-500';
    const iconClass = icon ? (iconPosition === 'left' ? 'pl-10' : 'pr-10') : '';

    return (
      <div className={`${widthClass} ${className}`}>
        {label && (
          <label htmlFor={inputId} className="block text-sm font-medium text-gray-700 mb-1">
            {label}
          </label>
        )}
        <div className="relative">
          {icon && iconPosition === 'left' && (
            <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-gray-500">
              {icon}
            </div>
          )}
          <input
            ref={ref}
            id={inputId}
            className={`appearance-none block px-3 py-2 border rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-2 ${errorClass} ${iconClass} ${widthClass}`}
            {...rest}
          />
          {icon && iconPosition === 'right' && (
            <div className="absolute inset-y-0 right-0 pr-3 flex items-center pointer-events-none text-gray-500">
              {icon}
            </div>
          )}
        </div>
        {error && <ErrorMessage message={error} />}
        {helperText && !error && (
          <p className="mt-1 text-sm text-gray-500">{helperText}</p>
        )}
      </div>
    );
  }
);

InputField.displayName = 'InputField';

export default InputField; 