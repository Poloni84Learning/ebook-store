const messages = {
  en: {
    home: {
      title: "Home",
      heroTitle: "Welcome to the EBookstore!",
      heroSubtitle: "Find the perfect book for every moment.",
      browseBooks: "Browse Books",
      heroImageAlt: "Bookshelf and reading illustration",
      featuredBooks: "Featured Books",
      ctaTitle: "Join our book community today!",
      ctaSubtitle: "Sign up and get access to exclusive deals, book recommendations, and more.",
      signUpNow: "Sign Up Now"
    },
    
    categories: 'Categories',
    login: 'Login',
    about: 'About Us',
    contact: 'Contact',
    faq: 'FAQ',
    shipping: 'Shipping',
    returns: 'Returns',
    footer: {
      description: 'Your trusted source for ebooks in every genre.',
      quickLinks: 'Quick Links',
      customerService: 'Customer Service',
      newsletter: 'Newsletter',
      subscribeText: 'Subscribe to receive updates, offers, and news.',
      emailPlaceholder: 'Enter your email',
      subscribe: 'Subscribe',
      rights: 'All rights reserved.',
    },
    "checkout": {
      "title": "Checkout",
      "steps": {
        "shipping": "Shipping Information",
        "payment": "Payment Details",
        "review": "Order Review"
      },
      "form": {
        "name": "Full Name",
        "email": "Email Address",
        "address": "Shipping Address",
        "city": "City",
        "zip": "ZIP Code"
      },
      "orderSummary": {
        "title": "Order Summary",
        "subtotal": "Subtotal",
        "shipping": "Shipping",
        "total": "Total",
        "checkoutButton": "Place Order"
      },
      "errors": {
        "required": "This field is required",
        "email": "Please enter a valid email"
      },
      "review": {
      "shippingTitle": "Shipping Information",
      "paymentTitle": "Payment Method",
      "orderTitle": "Your Order",
      "bankTransferNotice": "You will complete payment via bank transfer",
      "codNotice": "You will pay cash when your order arrives",
      "agreeTo": "I agree to the",
      "terms": "Terms and Conditions"
    },
  },

    
  orderHistory: {
    title: "Order History",
    emptyMessage: "You haven't placed any orders yet.",
    orderId: "Order ID",
    orderDate: "Order Date",
    itemsCount: "Items",
    customerInfo: "Customer Information",
    paymentMethod: "Payment Method",
    orderItems: "Order Items",
    total: "Total",
    download: "Download",
    downloadSuccess: "Download started successfully",
    downloadError: "Failed to download the book",
    status: {
      pending: "Pending",
      processing: "Processing",
      completed: "Completed",
      cancelled: "Cancelled"
    }
  },
      cart: {
    addedTitle: "Added to Cart",
    addedMessage: "The item has been added to your shopping cart.",
    continueShopping: "Continue Shopping",
    goToCart: "Go to Cart",
       title: "Your Cart ({count} items)",
    clearCart: "Clear Cart",
    empty: {
      title: "Your cart is empty",
      description: "Start shopping to add items to your cart"
    },
    subtotal: "Subtotal",
    shipping: "Shipping",
    shippingNote: "Shipping and taxes calculated at checkout",
    total: "Total",
    checkout: "Proceed to Checkout",
    
    clearModal: {
      title: "Clear your cart?",
      description: "This will remove all items from your cart. This action cannot be undone.",
      confirm: "Clear Cart"
    },
    deleteModal: {
      title: "Remove this item?",
      description: "This will remove the item from your cart. You can add it again later if you change your mind."
    },
    checkoutModal: {
      title: "Proceed to Checkout?",
      description: "Please review your order summary before proceeding:",
      confirm: "Confirm Checkout"
    }
  },
    "common": {
      "name": "Name",
      "email": "Email",
      "phone": "Phone",
      "zipCode": "ZIP Code",
      "quantity": "Quantity",
      "each": "each",
      "by":"by",
       cancel: "Cancel",
      submit: "Submit",
      submitting: "Submitting...",
      remove: "Remove",
      
    },
    books: {
    pageTitle: "Our Book Collection",
    pageSubtitle: "Discover amazing books from talented authors around the world",
    searchPlaceholder: "Search by title or author...",
    sortBy: "Sort by",
    newest: "Newest",
    priceLow: "Price: Low to High",
    priceHigh: "Price: High to Low",
    topRated: "Top Rated",
    noBooksFound: "No books found",
    tryDifferentSearch: "Try a different search term or filter",
    
    // Thêm các key khác nếu cần
  }  ,
  book:{
    by:"by",
        addToCart: "Add to Cart",
    description: "Description",
    language: "Language",
    isbn: "ISBN",
    category: "Category",
    relatedBooks: "You May Also Like",
    reviews: "Customer Reviews",
    addReview: "Add Review",
    writeReview: "Write a Review",
    rating: "Rating",
    comment: "Comment",
    reviewPlaceholder: "Share your thoughts about this book...",
    noReviews: "No reviews yet. Be the first to review!",
    alreadyReviewed: "You have already reviewed this book",
    reviewError: "Error submitting review. Please try again.",
    stock: "Stock",
    inStock: "In Stock",
    outOfStock: "Out of Stock",
        author: "Author",
    price: "Price",
    each: "each"
  }
    
  },
  ru: {
    home: {
      title: "Главная",
      heroTitle: "Добро пожаловать в EBookstore!",
      heroSubtitle: "Найдите идеальную книгу для любого момента.",
      browseBooks: "Смотреть книги",
      heroImageAlt: "Иллюстрация с полкой и чтением",
      featuredBooks: "Избранные книги",
      ctaTitle: "Присоединяйтесь к нашему книжному сообществу!",
      ctaSubtitle: "Зарегистрируйтесь и получите доступ к эксклюзивным предложениям, рекомендациям и многому другому.",
      signUpNow: "Зарегистрироваться"
    },
    books: 'Книги',
    categories: 'Категории',
    login: 'Войти',
    about: 'О нас',
    contact: 'Контакты',
    faq: 'Вопросы и ответы',
    shipping: 'Доставка',
    returns: 'Возвраты',
    footer: {
      description: 'Ваш надежный источник электронных книг в любом жанре.',
      quickLinks: 'Быстрые ссылки',
      customerService: 'Служба поддержки',
      newsletter: 'Рассылка',
      subscribeText: 'Подпишитесь на обновления, предложения и новости.',
      emailPlaceholder: 'Введите ваш email',
      subscribe: 'Подписаться',
      rights: 'Все права защищены.',
    }
  }
}


export default messages
